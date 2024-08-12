package db

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
)

// NOTE: // I would not put any logger here at all. The responsibility for logging should either be in the handlers of the server,
// or inside the service layer for certain scenarios. But most often you would place it in the handlers of the server (server/handlers_notes.go)
// logger is the interface that wraps around methods Debug, Info and Error.
/* type logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
} */

// cosmosClient is the interface that wraps around the Azure SDK CosmosDB client.
// NOTE: This can be renamed to client after the rest is implemented.
// NOTE: In this case id is a better name than itemID. Since we always know that the ID refers
// to the item ID in the CosmosDB.
// Avoid using prefixes or suffixes where it can be easily inferred by the context.
// Since in your scenario, the only thing we care about is the Value of the response, we can just return that.
// There might be cases where you want to get the status code or other information from the response. Then
// adjust accordingly.
type client interface {
	CreateItem(ctx context.Context, partitionKey string, item []byte) ([]byte, error)
	ReplaceItem(ctx context.Context, partitionKey, id string, item []byte) ([]byte, error)
	DeleteItem(ctx context.Context, partitionKey, id string) error
	ReadItem(ctx context.Context, partitionKey, id string) ([]byte, error)
	// Ideally you would add another parameter here. Like options or filter
	// to make the ListItems method more flexible, but that is out of scope for this
	// example.
	ListItems(ctx context.Context, partitionKey string) ([][]byte, error)
}

// NOTE: cosmosClient wraps around the "real" Azure SDK CosmosDB client.
type cosmosClient struct {
	*azcosmos.ContainerClient
}

// NewCosmosClient returns a new CosmosDB client. All it does is doing the needed setup for
// the components of the Azure SDK that we need.
func NewCosmosClient(connectionString, database, container string) (*cosmosClient, error) {
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
	}
	containerClient, err := client.NewContainer(database, container)
	if err != nil {
		return nil, err
	}
	return &cosmosClient{
		ContainerClient: containerClient,
	}, nil
}

// NOTE: This will not be tested by you, since the responsibility of testing the inner functionality
// is on the Azure SDK. You will only test the methods that use this client.
func (c cosmosClient) CreateItem(ctx context.Context, partitionKey string, item []byte) ([]byte, error) {
	resp, err := c.ContainerClient.CreateItem(ctx, azcosmos.NewPartitionKeyString(partitionKey), item, &azcosmos.ItemOptions{
		EnableContentResponseOnWrite: true,
	})
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// NOTE: This will not be tested by you, since the responsibility of testing the inner functionality
// is on the Azure SDK. You will only test the methods that use this client.
func (c cosmosClient) ReplaceItem(ctx context.Context, partitionKey string, id string, item []byte) ([]byte, error) {
	resp, err := c.ContainerClient.ReplaceItem(ctx, azcosmos.NewPartitionKeyString(partitionKey), id, item, &azcosmos.ItemOptions{
		EnableContentResponseOnWrite: true,
	})
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// NOTE: This will not be tested by you, since the responsibility of testing the inner functionality
// is on the Azure SDK. You will only test the methods that use this client.
func (c cosmosClient) DeleteItem(ctx context.Context, partitionKey string, id string) error {
	_, err := c.ContainerClient.DeleteItem(ctx, azcosmos.NewPartitionKeyString(partitionKey), id, nil)
	return err
}

// NOTE: This will not be tested by you, since the responsibility of testing the inner functionality
// is on the Azure SDK. You will only test the methods that use this client.
func (c cosmosClient) ReadItem(ctx context.Context, partitionKey string, id string) ([]byte, error) {
	resp, err := c.ContainerClient.ReadItem(ctx, azcosmos.NewPartitionKeyString(partitionKey), id, nil)
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// NOTE: This will not be tested by you, since the responsibility of testing the inner functionality
// is on the Azure SDK. You will only test the methods that use this client.
func (c cosmosClient) ListItems(ctx context.Context, partitionKey string) ([][]byte, error) {
	// As mentioned with the interface, this method should have options or filters
	// that builds the query according to the provided filter. So the list
	// method can be used for different scenarios.
	query := "SELECT * FROM c"

	var items [][]byte
	pager := c.ContainerClient.NewQueryItemsPager(query, azcosmos.NewPartitionKeyString(partitionKey), nil)
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		items = append(items, resp.Items...)
	}
	return items, nil
}

type CosmosDB struct {
	cl client
}

func NewCosmosDB(client client) (*CosmosDB, error) {
	if client == nil {
		return nil, ErrClientRequired
	}

	return &CosmosDB{
		cl: client,
	}, nil
}

var newUUID = func() string {
	return uuid.NewString()
}

// var newUUID func()string = uuid.NewString

func (c *CosmosDB) CreateNote(ctx context.Context, note Note) (Note, error) {
	note.ID = newUUID()

	bytes, err := json.Marshal(&note)
	if err != nil {
		return Note{}, err
	}

	data, err := c.cl.CreateItem(ctx, note.Category, bytes)
	if err != nil {
		return Note{}, checkError(err)
	}

	var n Note
	if err := json.Unmarshal(data, &n); err != nil {
		return Note{}, err
	}

	return n, nil
}

func (c *CosmosDB) UpdateNote(ctx context.Context, note Note) (Note, error) {
	bytes, err := json.Marshal(&note)
	if err != nil {
		return Note{}, err
	}

	data, err := c.cl.ReplaceItem(ctx, note.Category, note.ID, bytes)
	if err != nil {
		return Note{}, checkError(err)
	}

	var n Note
	if err := json.Unmarshal(data, &n); err != nil {
		return Note{}, err
	}

	return n, nil
}

func (c *CosmosDB) DeleteNote(ctx context.Context, id, category string) error {
	if err := c.cl.DeleteItem(ctx, category, id); err != nil {
		return checkError(err)
	}
	return nil
}

func (c *CosmosDB) GetNotesByCategory(ctx context.Context, category string) ([]Note, error) {
	data, err := c.cl.ListItems(ctx, category)
	if err != nil {
		return []Note{}, err
	}

	notes := make([]Note, len(data))
	for i, d := range data {
		if err := json.Unmarshal(d, &notes[i]); err != nil {
			return []Note{}, err
		}
	}
	return notes, nil
}

func (c *CosmosDB) GetNoteByID(ctx context.Context, category, id string) (Note, error) {
	data, err := c.cl.ReadItem(ctx, category, id)
	if err != nil {
		return Note{}, checkError(err)
	}

	var note Note
	if err = json.Unmarshal(data, &note); err != nil {
		return Note{}, err
	}

	return note, nil
}
