package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

// logger is the interface that wraps around methods Debug, Info and Error.
type logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type QueryResult struct {
}
type cosmosClient interface {
	CreateItem(ctx context.Context, partitionKey azcosmos.PartitionKey, item []byte, o *azcosmos.ItemOptions) ([]byte, error)
	ReplaceItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, item []byte, o *azcosmos.ItemOptions) ([]byte, error)
	DeleteItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, o *azcosmos.ItemOptions) ([]byte, error)
	ReadItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, o *azcosmos.ItemOptions) ([]byte, error)
	Query(ctx context.Context, query string, partitionKey azcosmos.PartitionKey, o *azcosmos.QueryOptions) ([][]byte, error)
}

type CosmosContainerClient struct {
	cl *azcosmos.ContainerClient
}

type NotesDB struct {
	cl  cosmosClient
	log logger
}

func NewNotesDB(client cosmosClient, logger logger) (*NotesDB, error) {
	if client == nil {
		return nil, ErrClientRequired
	}
	if logger == nil {
		return nil, ErrLoggerRequired
	}

	return &NotesDB{
		cl:  client,
		log: logger,
	}, nil
}

var newUUID = func() string {
	return uuid.NewString()
}

// var newUUID func()string = uuid.NewString

func (c *NotesDB) CreateNote(ctx context.Context, note Note) (Note, error) {
	// assign Note ID if it is not set
	if len(note.ID) == 0 {
		note.ID = newUUID()
	}

	// assign current time if CreatedAt is not set
	if note.CreatedAt.IsZero() {
		note.CreatedAt = time.Now().UTC()
	}

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
	if err := json.Unmarshal(resp, &noteDB); err != nil {
		c.log.Error("Failed unmarshal the note.", logError(err)...)
		return Note{}, err
	}
	return noteDB, nil
}

func (c *NotesDB) UpdateNote(ctx context.Context, note Note) (Note, error) {
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

	if err := json.Unmarshal(resp, &note); err != nil {
		c.log.Error("Failed unmarshal the note.", logError(err)...)
		return Note{}, err
	}
	return note, nil
}

func (c *NotesDB) DeleteNote(ctx context.Context, id, category string) error {
	pk := azcosmos.NewPartitionKeyString(category)

	// Q: would it a better practice to write a custom error message here, i.e. "Failed to delete a note in CosmosDB"?
	if _, err := c.cl.DeleteItem(ctx, pk, id, nil); err != nil {
		return checkError(err)
	}
	return nil
}

func (c *NotesDB) GetNotesByCategory(ctx context.Context, category string) ([]Note, error) {
	var notes []Note
	query := "SELECT * FROM c"
	pk := azcosmos.NewPartitionKeyString(category)
	respItems, err := c.cl.Query(ctx, query, pk, nil)
	if err != nil {
		return []Note{}, err
	}
	for _, item := range respItems {
		var note Note
		if err = json.Unmarshal(item, &note); err != nil {
			return []Note{}, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (c *NotesDB) GetNoteByID(ctx context.Context, category, id string) (Note, error) {
	pk := azcosmos.NewPartitionKeyString(category)
	// read the item from the container
	response, err := c.cl.ReadItem(ctx, pk, id, nil)
	// Q: would it a better practice to write a custom error message here, i.e. "Failed to get a note from the CosmosDB"?
	if err != nil {
		return Note{}, checkError(err)
	}

	var note Note
	if err = json.Unmarshal(response, &note); err != nil {
		return Note{}, err
	}
	return note, nil
}

func NewCosmosContainerClient(connectionString, databaseID, containerID string) (cosmosClient, error) {
	cosmosClient, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
	}

	databaseClient, err := cosmosClient.NewDatabase(databaseID)
	if err != nil {
		return nil, err
	}

	containerClient, err := databaseClient.NewContainer(containerID)
	if err != nil {
		return nil, err
	}

	_, err = containerClient.Read(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrClientConnection, err)
	}

	return &CosmosContainerClient{
		cl: containerClient,
	}, nil
}

func (c *CosmosContainerClient) CreateItem(ctx context.Context, partitionKey azcosmos.PartitionKey, item []byte, o *azcosmos.ItemOptions) ([]byte, error) {
	resp, err := c.cl.CreateItem(ctx, partitionKey, item, o)
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

func (c *CosmosContainerClient) ReplaceItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, item []byte, o *azcosmos.ItemOptions) ([]byte, error) {
	resp, err := c.cl.ReplaceItem(ctx, partitionKey, itemId, item, o)
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

func (c *CosmosContainerClient) DeleteItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, o *azcosmos.ItemOptions) ([]byte, error) {
	resp, err := c.cl.DeleteItem(ctx, partitionKey, itemId, o)
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

func (c *CosmosContainerClient) ReadItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, o *azcosmos.ItemOptions) ([]byte, error) {
	resp, err := c.cl.ReadItem(ctx, partitionKey, itemId, o)
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

func (c *CosmosContainerClient) Query(ctx context.Context, query string, partitionKey azcosmos.PartitionKey, o *azcosmos.QueryOptions) ([][]byte, error) {
	pager := c.cl.NewQueryItemsPager(query, partitionKey, nil)
	var items [][]byte
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		items = append(items, resp.Items...)
	}
	return items, nil
}
