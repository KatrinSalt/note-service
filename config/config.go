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
	Log      Logger
}

type Note struct {
	Timeout time.Duration
}

type Database struct {
	ConnectionString string `env:"COSMOSDB_CONNECTION_STRING,required"`
	Database         string `env:"COSMOSDB_DATABASE,required"`
	Container        string `env:"COSMOSDB_CONTAINER,required"`
	Log              Logger
}

type Logger struct {
	ServiceLevel string `env:"SERVICE_LOG_LEVEL"`
	DBLevel      string `env:"DB_LOG_LEVEL"`
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
				Database:  defaultCosmosDatabase,
				Container: defaultCosmosContainer,
				Log: Logger{
					DBLevel: defaultDBLogLevel,
				},
			},
			Log: Logger{
				ServiceLevel: defaultServiceLogLevel,
			},
		},
	}

	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
