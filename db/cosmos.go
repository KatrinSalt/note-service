package db

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

// logger is the interface that wraps around methods Debug, Info and Error.
type logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type client interface {
	CreateItem(ctx context.Context, partitionKey azcosmos.PartitionKey, item []byte, o *azcosmos.ItemOptions) (azcosmos.ItemResponse, error)
	ReplaceItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, item []byte, o *azcosmos.ItemOptions) (azcosmos.ItemResponse, error)
	DeleteItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, o *azcosmos.ItemOptions) (azcosmos.ItemResponse, error)
	ReadItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, o *azcosmos.ItemOptions) (azcosmos.ItemResponse, error)
	NewQueryItemsPager(query string, partitionKey azcosmos.PartitionKey, o *azcosmos.QueryOptions) *runtime.Pager[azcosmos.QueryItemsResponse]
}

type CosmosDB struct {
	cl  client
	log logger
}

// func NewCosmosDB(connectionString, dbID, containerID string, logger logger, client Client) (*CosmosDB, error) {
// 	if len(connectionString) == 0 {
// 		return nil, ErrConnStringRequired
// 	}
// 	if len(dbID) == 0 {
// 		return nil, ErrDbIdRequired
// 	}
// 	if len(containerID) == 0 {
// 		return nil, ErrContainerIdRequired
// 	}
// 	if logger == nil {
// 		return nil, ErrLoggerRequired
// 	}

// 	logger.Debug("Creating a new CosmosDB client.", "dbID", dbID, "containerID", containerID)

// 	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	database, err := client.NewDatabase(dbID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	container, err := database.NewContainer(containerID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &CosmosDB{
// 		client:    client,
// 		database:  database,
// 		container: container,
// 		log:       logger,
// 	}, nil
// }

func NewCosmosDB(client client, logger logger) (*CosmosDB, error) {
	if client == nil {
		return nil, ErrClientRequired
	}
	if logger == nil {
		return nil, ErrLoggerRequired
	}

	return &CosmosDB{
		cl:  client,
		log: logger,
	}, nil
}

var newUUID = func() string {
	return uuid.NewString()
}

// var newUUID func()string = uuid.NewString

func (c *CosmosDB) CreateNote(ctx context.Context, note Note) (Note, error) {
	note.ID = newUUID()

	// Q: would you 'properly' handle the error here, i.e. with the specific message?
	bytes, err := json.Marshal(&note)
	if err != nil {
		c.log.Error("Failed to marshal the note.", logError(err)...)
		return Note{}, err
	}

	pk := azcosmos.NewPartitionKeyString(note.Category)

	// Q: would it a better practice to write a custom error message here, i.e. "Failed to create a note in CosmosDB"?
	resp, err := c.cl.CreateItem(ctx, pk, bytes, &azcosmos.ItemOptions{
		EnableContentResponseOnWrite: true,
	})
	if err != nil {
		c.log.Error("Failed to create the note.", logError(err)...)
		return Note{}, checkError(err)
	}

	var noteDB Note
	if err := json.Unmarshal(resp.Value, &noteDB); err != nil {
		c.log.Error("Failed unmarshal the note.", logError(err)...)
		return Note{}, err
	}

	return noteDB, nil
}

func (c *CosmosDB) UpdateNote(ctx context.Context, note Note) (Note, error) {
	// Q: would you 'properly' handle the error here, i.e. with the specific message?
	bytes, err := json.Marshal(&note)
	if err != nil {
		c.log.Error("Failed to marshal the note.", logError(err)...)
		return Note{}, err
	}

	pk := azcosmos.NewPartitionKeyString(note.Category)

	// Q: would it a better practice to write a custom error message here, i.e. "Failed to update a note in CosmosDB"?
	resp, err := c.cl.ReplaceItem(ctx, pk, note.ID, bytes, &azcosmos.ItemOptions{
		EnableContentResponseOnWrite: true,
	})
	if err != nil {
		c.log.Error("Failed to update the note.", logError(err)...)
		return Note{}, checkError(err)
	}

	if err := json.Unmarshal(resp.Value, &note); err != nil {
		c.log.Error("Failed unmarshal the note.", logError(err)...)
		return Note{}, err
	}

	return note, nil
}

func (c *CosmosDB) DeleteNote(ctx context.Context, id, category string) error {
	pk := azcosmos.NewPartitionKeyString(category)

	// Q: would it a better practice to write a custom error message here, i.e. "Failed to delete a note in CosmosDB"?
	if _, err := c.cl.DeleteItem(ctx, pk, id, nil); err != nil {
		return checkError(err)
	}
	return nil
}

func (c *CosmosDB) GetNotesByCategory(ctx context.Context, category string) ([]Note, error) {
	var notes []Note
	query := "SELECT * FROM c"
	pk := azcosmos.NewPartitionKeyString(category)
	pager := c.cl.NewQueryItemsPager(query, pk, nil)
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
	response, err := c.cl.ReadItem(ctx, pk, id, nil)
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
