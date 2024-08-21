package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

type client interface {
	CreateItem(ctx context.Context, partitionKey string, item []byte) ([]byte, error)
	ReplaceItem(ctx context.Context, partitionKey string, id string, item []byte) ([]byte, error)
	DeleteItem(ctx context.Context, partitionKey string, id string) error
	ReadItem(ctx context.Context, partitionKey string, id string) ([]byte, error)
	ListItems(ctx context.Context, partitionKey string) ([][]byte, error)
}

type CosmosContainerClient struct {
	cl *azcosmos.ContainerClient
}

func NewCosmosContainerClient(connectionString, databaseID, containerID string) (*CosmosContainerClient, error) {
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
	}

	containerClient, err := client.NewContainer(databaseID, containerID)
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

func (c *CosmosContainerClient) CreateItem(ctx context.Context, partitionKey string, item []byte) ([]byte, error) {
	resp, err := c.cl.CreateItem(ctx, azcosmos.NewPartitionKeyString(partitionKey), item, &azcosmos.ItemOptions{
		EnableContentResponseOnWrite: true,
	})
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

func (c *CosmosContainerClient) ReplaceItem(ctx context.Context, partitionKey string, id string, item []byte) ([]byte, error) {
	resp, err := c.cl.ReplaceItem(ctx, azcosmos.NewPartitionKeyString(partitionKey), id, item, &azcosmos.ItemOptions{
		EnableContentResponseOnWrite: true,
	})
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

func (c *CosmosContainerClient) DeleteItem(ctx context.Context, partitionKey string, id string) error {
	_, err := c.cl.DeleteItem(ctx, azcosmos.NewPartitionKeyString(partitionKey), id, &azcosmos.ItemOptions{
		EnableContentResponseOnWrite: true,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *CosmosContainerClient) ReadItem(ctx context.Context, partitionKey string, id string) ([]byte, error) {
	resp, err := c.cl.ReadItem(ctx, azcosmos.NewPartitionKeyString(partitionKey), id, &azcosmos.ItemOptions{
		EnableContentResponseOnWrite: true,
	})
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

func (c *CosmosContainerClient) ListItems(ctx context.Context, partitionKey string) ([][]byte, error) {
	query := "SELECT * FROM c"
	pager := c.cl.NewQueryItemsPager(query, azcosmos.NewPartitionKeyString(partitionKey), nil)
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

type NotesDB struct {
	cl client
}

func NewNotesDB(client client) (*NotesDB, error) {
	if client == nil {
		return nil, ErrClientRequired
	}

	return &NotesDB{
		cl: client,
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
		return Note{}, err
	}

	// Q: would it a better practice to write a custom error message here, i.e. "Failed to create a note in CosmosDB"?
	resp, err := c.cl.CreateItem(ctx, note.Category, bytes)
	if err != nil {
		return Note{}, checkError(err)
	}

	var noteDB Note
	if err := json.Unmarshal(resp, &noteDB); err != nil {
		return Note{}, err
	}
	return noteDB, nil
}

func (c *NotesDB) UpdateNote(ctx context.Context, note Note) (Note, error) {
	bytes, err := json.Marshal(&note)
	if err != nil {
		return Note{}, err
	}

	// Q: would it a better practice to write a custom error message here, i.e. "Failed to update a note in CosmosDB"?
	resp, err := c.cl.ReplaceItem(ctx, note.Category, note.ID, bytes)
	if err != nil {
		return Note{}, checkError(err)
	}

	var noteDB Note
	if err := json.Unmarshal(resp, &noteDB); err != nil {
		return Note{}, err
	}
	return noteDB, nil
}

func (c *NotesDB) DeleteNote(ctx context.Context, id, category string) error {
	err := c.cl.DeleteItem(ctx, category, id)
	if err != nil {
		return checkError(err)
	}
	return nil
}

func (c *NotesDB) GetNotesByCategory(ctx context.Context, category string) ([]Note, error) {
	var notes []Note
	respItems, err := c.cl.ListItems(ctx, category)
	if err != nil {
		return []Note{}, checkError(err)
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
	// read the item from the container
	response, err := c.cl.ReadItem(ctx, category, id)
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
