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
	logger, err := setupLogger(config.Log.ServiceLevel)
	if err != nil {
		return nil, err
	}

	notesDB, err := setupNotesDB(config.Database)
	if err != nil {
		return nil, err
	}

	notesvc, err := notes.NewService(notesDB, logger, func(o *notes.ServiceOptions) {
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

func setupNotesDB(config Database) (*db.NotesDB, error) {
	if len(config.CosmosContainerClient.ConnectionString) == 0 {
		return nil, errors.New("cosmosdb connection string is empty")
	}
	if len(config.CosmosContainerClient.DatabaseID) == 0 {
		return nil, errors.New("cosmosdb database id is empty")
	}
	if len(config.CosmosContainerClient.DatabaseID) == 0 {
		return nil, errors.New("cosmosdb container id is empty")
	}

	containerClient, err := db.NewCosmosContainerClient(config.CosmosContainerClient.ConnectionString,
		config.CosmosContainerClient.DatabaseID, config.CosmosContainerClient.ContainerID)
	if err != nil {
		return nil, err
	}

	notesDB, err := db.NewNotesDB(containerClient)
	if err != nil {
		return nil, err
	}
	return notesDB, nil
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
