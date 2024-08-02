package config

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/KatrinSalt/notes-service/db"
	"github.com/KatrinSalt/notes-service/log"
	"github.com/KatrinSalt/notes-service/notes"
)

type services struct {
	Note notes.Service
}

func SetupServices(config Services) (*services, error) {
	logger, err := setupLogger(config.Log.ServiceLevel)
	if err != nil {
		return nil, err
	}

	cosmosDB, err := setupCosmosDB(config.Database)
	if err != nil {
		return nil, err
	}

	notesvc, err := notes.NewService(cosmosDB, logger, func(o *notes.ServiceOptions) {
		o.Timeout = config.Note.Timeout
	},
	)

	if err != nil {
		return nil, err
	}

	return &services{
		Note: notesvc,
	}, nil

}

func setupCosmosContainerClient(config Client) (*azcosmos.ContainerClient, error) {
	if len(config.ConnectionString) == 0 {
		return nil, errors.New("cosmosdb connection string is empty")
	}
	if len(config.DatabaseID) == 0 {
		return nil, errors.New("cosmosdb database id is empty")
	}
	if len(config.ContainerID) == 0 {
		return nil, errors.New("cosmosdb container id is empty")
	}

	CosmosContainerClient, err := azcosmos.NewClientFromConnectionString(config.ConnectionString, nil)
	if err != nil {
		return nil, err
	}

	databaseClient, err := CosmosContainerClient.NewDatabase(config.DatabaseID)
	if err != nil {
		return nil, err
	}

	containerClient, err := databaseClient.NewContainer(config.ContainerID)
	if err != nil {
		return nil, err
	}

	return containerClient, nil
}

func setupCosmosDB(config Database) (*db.CosmosDB, error) {
	client, err := setupCosmosContainerClient(config.CosmosContainerClient)
	if err != nil {
		return nil, err
	}

	logger, err := setupLogger(config.Log.DBLevel)
	if err != nil {
		return nil, err
	}

	cosmosDB, err := db.NewCosmosDB(client, logger)
	if err != nil {
		return nil, err
	}

	return cosmosDB, nil
}

func setupLogger(logLevel string) (*log.Logger, error) {
	if len(logLevel) == 0 {
		return log.New(), nil
	} else {
		logger, err := log.NewWithSetLevel(logLevel)
		if err != nil {
			return nil, err
		}
		return logger, nil
	}
}
