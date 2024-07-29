package config

import (
	"errors"

	"github.com/KatrinSalt/notes-service/db"
	"github.com/KatrinSalt/notes-service/log"
	"github.com/KatrinSalt/notes-service/notes"
)

type services struct {
	Note notes.Service
}

func SetupServices(config Services) (*services, error) {
	logger, err := setupServiceLogger(config.Log)
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

func setupCosmosDB(config Database) (*db.CosmosDB, error) {
	if len(config.Cosmos.ConnectionString) == 0 {
		return nil, errors.New("cosmosdb connection string is empty")
	} else if len(config.Cosmos.DatabaseID) == 0 {
		return nil, errors.New("cosmosdb database id is empty")
	} else if len(config.Cosmos.ContainerID) == 0 {
		return nil, errors.New("cosmosdb container id is empty")
	} else {
		logger, err := setupDBLogger(config.Log)
		if err != nil {
			return nil, err
		}
		cosmosDB, err := db.NewCosmosDB(config.Cosmos.ConnectionString, config.Cosmos.DatabaseID, config.Cosmos.ContainerID, logger)
		if err != nil {
			return nil, err
		}
		return cosmosDB, nil
	}
}

func setupServiceLogger(config Logger) (*log.Logger, error) {
	if len(config.ServiceLevel) == 0 {
		return log.New(), nil
	} else {
		logger, err := log.NewWithSetLevel(config.ServiceLevel)
		if err != nil {
			return nil, err
		}
		return logger, nil
	}
}

func setupDBLogger(config Logger) (*log.Logger, error) {
	if len(config.DBLevel) == 0 {
		return log.New(), nil
	} else {
		logger, err := log.NewWithSetLevel(config.DBLevel)
		if err != nil {
			return nil, err
		}
		return logger, nil
	}
}
