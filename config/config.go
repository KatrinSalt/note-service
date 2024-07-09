package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

// Configuration contains the configuration for the application.
type Configuration struct {
	Server   Server
	Services Services
}

// Server contains the configuration for the server.
type Server struct {
	Host string
	Port string
}

type Services struct {
	Note     Note
	Database Database
}

type Note struct {
	Timeout time.Duration
}

type Database struct {
	Cosmos DatabaseCosmos
}

type DatabaseCosmos struct {
	ConnectionString string `env:"COSMOSDB_CONNECTION_STRING,required"`
	DatabaseID       string `env:"COSMOSDB_DATABASE_ID"`
	ContainerID      string `env:"COSMOSDB_CONTAINER_ID"`
}

// Options for the configuration.
type Options struct{}

// Option is a function that sets options for the configuration.
type Option func(o *Options)

// New creates a new configuration.
func New(options ...Option) (Configuration, error) {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}

	cfg := Configuration{
		Server: Server{
			Host: defaultServerHost,
			Port: defaultServerPort,
		},
		Services: Services{
			Note: Note{
				Timeout: defaultNoteTimeout,
			},
			Database: Database{
				Cosmos: DatabaseCosmos{
					DatabaseID:  defaultCosmosDatabaseID,
					ContainerID: defaultCosmosContainerID,
				},
			},
		},
	}

	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
