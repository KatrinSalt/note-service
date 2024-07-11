package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

const (
	counterContainerID = "counter"
	counterItemID      = "1"
)

var pkCounter = azcosmos.NewPartitionKeyString(counterItemID)

type Database interface {
	// CreateNote creates a new note.
	CreateNote(ctx context.Context, note *Note) error
}

type cosmosDB struct {
	client    *azcosmos.Client
	database  *azcosmos.DatabaseClient
	container *azcosmos.ContainerClient
	counter   *azcosmos.ContainerClient
}

func NewCosmosDB(connectionString, dbID, containerID string) (*cosmosDB, error) {
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

	counter, err := database.NewContainer(counterContainerID)
	if err != nil {
		return nil, err
	}

	return &cosmosDB{
		client:    client,
		database:  database,
		container: container,
		counter:   counter,
	}, nil
}

func (c *cosmosDB) CreateNote(ctx context.Context, note *Note) error {
	if err := c.assignID(ctx, note); err != nil {
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

	// update the counter item with the incremented ID
	id, err := strconv.Atoi(note.ID)
	if err != nil {
		fmt.Printf("Failed to convert the ID to an integer: %s\n", err)
	}

	counter := Counter{
		ID:    counterItemID,
		MaxID: id,
	}
	if err := c.updateCounterItem(ctx, counter); err != nil {
		fmt.Printf("Failed to update the counter item: %s\n", err)
	}

	return nil
}

func (c *cosmosDB) assignID(ctx context.Context, note *Note) error {
	id, err := c.generateID(ctx)
	if err != nil {
		fmt.Printf("Failed to generate an ID: %s\n", err)
		return err
	}

	fmt.Printf("Generated ID: %s\n", id)

	note.ID = id
	fmt.Printf("Note struct with ID: %+v\n", note)
	return nil
}

func (c *cosmosDB) generateID(ctx context.Context) (string, error) {
	var id string
	counter := Counter{
		ID:    counterItemID,
		MaxID: 0,
	}

	// check if the counter item exists in counter container
	exists, err := c.checkIfCounterExists(ctx, counter)
	if err != nil {
		return id, err
	}

	// if the counter item does not exist, create it
	if !exists {
		if err := c.instantiateCounter(ctx, counter); err != nil {
			return id, err
		}
	}

	// get the recent counter item stored in counter container
	itemResponse, err := c.counter.ReadItem(ctx, pkCounter, counterItemID, nil)
	if err != nil {
		fmt.Printf("Failed to read the counter item: %s\n", err)
		return id, checkError(err)
	}
	if err = json.Unmarshal(itemResponse.Value, &counter); err != nil {
		fmt.Printf("Failed to unmarshal the counter: %s\n", err)
		return id, err
	}

	// increment the counter
	counter.MaxID++
	fmt.Printf("Counter struct: %+v\n", counter)

	id = strconv.Itoa(counter.MaxID)
	return id, nil
}

func (c *cosmosDB) updateCounterItem(ctx context.Context, counter Counter) error {
	bytes, err := json.Marshal(counter)
	if err != nil {
		fmt.Printf("Failed to marshal the counter item: %s\n", err)
		return err
	}

	_, err = c.counter.ReplaceItem(ctx, pkCounter, counterItemID, bytes, nil)
	if err != nil {
		fmt.Printf("Failed to update the counter item: %s\n", err)
		return checkError(err)
	}
	return nil
}

func (c *cosmosDB) checkIfCounterExists(ctx context.Context, counter Counter) (bool, error) {
	query := "SELECT VALUE COUNT(1) FROM c"
	pk := azcosmos.NewPartitionKeyString(counter.ID)
	queryPager := c.counter.NewQueryItemsPager(query, pk, nil)
	for queryPager.More() {
		page, err := queryPager.NextPage(ctx)
		if err != nil {
			fmt.Printf("Failed to queue NoteDB counter container: %s\n", err)
			return false, err
		}
		for _, item := range page.Items {
			var count int
			if err := json.Unmarshal(item, &count); err != nil {
				fmt.Printf("Failed to unmarshal the counter: %s\n", err)
				return false, err
			}
			if count == 0 {
				return false, nil
			}
		}
	}
	return true, nil
}

func (c *cosmosDB) instantiateCounter(ctx context.Context, counter Counter) error {
	bytes, err := json.Marshal(counter)
	if err != nil {
		fmt.Printf("Failed to marshal the counter item: %s\n", err)
		return err
	}

	pk := azcosmos.NewPartitionKeyString(counter.ID)

	if _, err := c.counter.CreateItem(ctx, pk, bytes, nil); err != nil {
		fmt.Printf("Failed to create a counter item in NoteDB counter container: %s\n", err)
		return checkError(err)
	}
	return nil
}
