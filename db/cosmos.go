package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Database interface {
	// CreateNote creates a new note.
	CreateNote(ctx context.Context, note Note) error
}

type cosmosDB struct {
	client    *azcosmos.Client
	database  *azcosmos.DatabaseClient
	container *azcosmos.ContainerClient
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

	// FYI: temporary for debugging
	fmt.Printf("Get container:\t%s\n", container.ID())

	return &cosmosDB{
		client:    client,
		database:  database,
		container: container,
	}, nil
}

func (c *cosmosDB) CreateNote(ctx context.Context, note Note) error {
	bytes, err := json.Marshal(note)
	if err != nil {
		fmt.Printf("Failed to marshal the note: %s\n", err)
		return err
	}

	if _, err := c.container.CreateItem(ctx, azcosmos.NewPartitionKeyString(note.ID), bytes, nil); err != nil {
		fmt.Printf("Failed to create a note in CosmosDB: %s\n", err)
		return checkError(err)
	}

	return nil
}
