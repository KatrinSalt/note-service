package db

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

// logger is the interface that wraps around methods Debug, Info and Error.
type logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type CosmosDB struct {
	client    *azcosmos.Client
	database  *azcosmos.DatabaseClient
	container *azcosmos.ContainerClient
	log       logger
}

func NewCosmosDB(connectionString, dbID, containerID string, logger logger) (*CosmosDB, error) {
	if len(connectionString) == 0 {
		return nil, ErrConnStringRequired
	}
	if len(dbID) == 0 {
		return nil, ErrDbIdRequired
	}
	if len(containerID) == 0 {
		return nil, ErrContainerIdRequired
	}
	if logger == nil {
		return nil, ErrLoggerRequired
	}

	logger.Debug("Creating a new CosmosDB client.", "dbID", dbID, "containerID", containerID)

	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
	}

	database, err := client.NewDatabase(dbID)
	if err != nil {
		return nil, err
	}

	container, err := database.NewContainer(containerID)
	if err != nil {
		return nil, err
	}

	return &CosmosDB{
		client:    client,
		database:  database,
		container: container,
		log:       logger,
	}, nil
}

var newUUID = func() string {
	return uuid.NewString()
}

func (c *CosmosDB) CreateNote(ctx context.Context, note *Note) error {
	note.ID = newUUID()

	// Q: would you 'properly' handle the error here, i.e. with the specific message?
	bytes, err := json.Marshal(&note)
	if err != nil {
		return err
	}

	pk := azcosmos.NewPartitionKeyString(note.Category)

	// Q: would it a better practice to write a custom error message here, i.e. "Failed to create a note in CosmosDB"?
	if _, err := c.container.CreateItem(ctx, pk, bytes, nil); err != nil {
		return checkError(err)
	}
	return nil
}

func (c *CosmosDB) UpdateNote(ctx context.Context, note *Note) error {
	// Q: would you 'properly' handle the error here, i.e. with the specific message?
	bytes, err := json.Marshal(&note)
	if err != nil {
		return err
	}

	pk := azcosmos.NewPartitionKeyString(note.Category)

	// Q: would it a better practice to write a custom error message here, i.e. "Failed to update a note in CosmosDB"?
	if _, err := c.container.ReplaceItem(ctx, pk, note.ID, bytes, nil); err != nil {
		return checkError(err)
	}
	return nil
}

func (c *CosmosDB) DeleteNote(ctx context.Context, id, category string) error {
	pk := azcosmos.NewPartitionKeyString(category)

	// Q: would it a better practice to write a custom error message here, i.e. "Failed to delete a note in CosmosDB"?
	if _, err := c.container.DeleteItem(ctx, pk, id, nil); err != nil {
		return checkError(err)
	}
	return nil
}

func (c *CosmosDB) GetNotesByCategory(ctx context.Context, category string) ([]Note, error) {
	var notes []Note
	query := "SELECT * FROM c"
	pk := azcosmos.NewPartitionKeyString(category)
	pager := c.container.NewQueryItemsPager(query, pk, nil)
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return []Note{}, err
		}

		for _, item := range resp.Items {
			var note Note
			if err = json.Unmarshal(item, &note); err != nil {
				return []Note{}, err
			}
			notes = append(notes, note)
		}
	}
	return notes, nil
}

func (c *CosmosDB) GetNoteByID(ctx context.Context, category, id string) (Note, error) {
	pk := azcosmos.NewPartitionKeyString(category)
	// read the item from the container
	response, err := c.container.ReadItem(ctx, pk, id, nil)
	// Q: would it a better practice to write a custom error message here, i.e. "Failed to get a note from the CosmosDB"?
	if err != nil {
		return Note{}, checkError(err)
	}

	var note Note
	if err = json.Unmarshal(response.Value, &note); err != nil {
		return Note{}, err
	}

	return note, nil
}
