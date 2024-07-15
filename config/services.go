package config

import (
	"errors"

	"github.com/KatrinSalt/notes-service/db"
	"github.com/KatrinSalt/notes-service/notes"
)

type services struct {
	Note notes.Service
}

func SetupServices(config Services) (*services, error) {
	cosmosDB, err := setupCosmosDB(config.Database)
	if err != nil {
		return nil, err
	}

	notesvc, err := notes.NewService(cosmosDB, func(o *notes.ServiceOptions) {
		o.Timeout = config.Note.Timeout
	})
	if err != nil {
		return nil, err
	}

	return &services{
		Note: notesvc,
	}, nil

}

func setupCosmosDB(config Database) (*db.CosmosDB, error) {
	if len(config.Cosmos.ConnectionString) == 0 {
		return nil, errors.New("cosmosdb connection string is empty")
	} else if len(config.Cosmos.DatabaseID) == 0 {
		return nil, errors.New("cosmosdb database id is empty")
	} else if len(config.Cosmos.ContainerID) == 0 {
		return nil, errors.New("cosmosdb container id is empty")
	} else {
		cosmosDB, err := db.NewCosmosDB(config.Cosmos.ConnectionString, config.Cosmos.DatabaseID, config.Cosmos.ContainerID)
		if err != nil {
			return nil, err
		}
		return cosmosDB, nil
	}
}
