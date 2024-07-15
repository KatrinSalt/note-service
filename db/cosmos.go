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

func (c *CosmosDB) assignID(note *Note) error {
	note.ID = uuid.New()

	fmt.Printf("Note struct with ID: %+v\n", note)
	return nil
}
