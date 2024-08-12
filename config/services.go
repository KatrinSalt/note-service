package config

import (
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

func setupCosmosDB(config Database) (*db.CosmosDB, error) {
	client, err := db.NewCosmosClient(config.ConnectionString, config.Database, config.Container)
	if err != nil {
		return nil, err
	}

	cosmosDB, err := db.NewCosmosDB(client)
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
