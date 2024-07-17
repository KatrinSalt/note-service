package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

type CosmosDB struct {
	client    *azcosmos.Client
	database  *azcosmos.DatabaseClient
	container *azcosmos.ContainerClient
}

func NewCosmosDB(connectionString, dbID, containerID string) (*CosmosDB, error) {
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
	}

	// FYI: temporary for debugging
	fmt.Println("Starting connection to CosmosDB...")

	database, err := client.NewDatabase(dbID)
	if err != nil {
		return nil, err
	}

	// FYI: temporary for debugging
	fmt.Printf("Get database:\t%s\n", database.ID())

	container, err := database.NewContainer(containerID)
	if err != nil {
		return nil, err
	}

	return &CosmosDB{
		client:    client,
		database:  database,
		container: container,
	}, nil
}

func (c *CosmosDB) CreateNote(ctx context.Context, note *Note) error {
	if err := c.assignID(note); err != nil {
		fmt.Printf("Failed to assign an ID: %s\n", err)
		return err
	}

	fmt.Printf("Note struct which is sent to DB: %+v\n", note)

	bytes, err := json.Marshal(&note)
	if err != nil {
		fmt.Printf("Failed to marshal the note: %s\n", err)
		return err
	}

	pk := azcosmos.NewPartitionKeyString(note.Category)

	if _, err := c.container.CreateItem(ctx, pk, bytes, nil); err != nil {
		fmt.Printf("Failed to create a note in CosmosDB: %s\n", err)
		return checkError(err)
	}
	return nil
}

func (c *CosmosDB) UpdateNote(ctx context.Context, note *Note) error {
	fmt.Printf("Note struct which is sent to DB: %+v\n", note)

	bytes, err := json.Marshal(&note)
	if err != nil {
		fmt.Printf("Failed to marshal the note: %s\n", err)
		return err
	}

	pk := azcosmos.NewPartitionKeyString(note.Category)

	if _, err := c.container.ReplaceItem(ctx, pk, note.ID.String(), bytes, nil); err != nil {
		fmt.Printf("Failed to update a note in CosmosDB: %s\n", err)
		return checkError(err)
	}
	return nil
}

func (c *CosmosDB) DeleteNote(ctx context.Context, id, category string) error {
	pk := azcosmos.NewPartitionKeyString(category)

	if _, err := c.container.DeleteItem(ctx, pk, id, nil); err != nil {
		fmt.Printf("Failed to delete a note in CosmosDB: %s\n", err)
		return checkError(err)
	}
	return nil
}

func (c *CosmosDB) GetNotesByCategory(ctx context.Context, category string) ([]Note, error) {
	var notes []Note
	query := "SELECT * FROM c"
	pk := azcosmos.NewPartitionKeyString(category)
	queryPager := c.container.NewQueryItemsPager(query, pk, nil)
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(ctx)
		if err != nil {
			fmt.Printf("Failed to query NoteDB note container: %s\n", err)
			return []Note{}, err
		}

		for _, item := range queryResponse.Items {
			var note Note
			if err = json.Unmarshal(item, &note); err != nil {
				fmt.Printf("Failed to unmarshal the container note response: %s\n", err)
				return []Note{}, err
			}
			notes = append(notes, note)
		}
	}
	fmt.Printf("Notes from DB in category %s: %+v\n", category, notes)
	return notes, nil
}

func (c *CosmosDB) GetNoteByID(ctx context.Context, category, id string) (Note, error) {
	pk := azcosmos.NewPartitionKeyString(category)
	// read the item from the container
	response, err := c.container.ReadItem(ctx, pk, id, nil)
	if err != nil {
		fmt.Printf("Failed to read a note from CosmosDB: %s\n", err)
		return Note{}, checkError(err)
	}

	var note Note
	if err = json.Unmarshal(response.Value, &note); err != nil {
		fmt.Printf("Failed to unmarshal the container note response: %s\n", err)
		return Note{}, err
	}

	return note, nil
}

func (c *CosmosDB) assignID(note *Note) error {
	note.ID = uuid.New()

	fmt.Printf("Note struct with ID: %+v\n", note)
	return nil
}
