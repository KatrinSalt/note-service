package config

import "time"

// Default server configuration.
const (
	defaultServerHost = "localhost"
	defaultServerPort = "3000"
)

// Default Note configuration.
const (
	defaultNoteTimeout = 10 * time.Second
)

// Default CosmosDB configuration.
const (
	// The application is called notes (or notes-service). Hence a fitting name
	// is notes.
	defaultCosmosDatabase = "notes"
	// The container can also be called notes, since it contains the notes.
	// if other containers are needed for other purposes, they can be named accordingly.
	defaultCosmosContainer = "notes"
)

// Default Logger configuration for Service.
const (
	defaultServiceLogLevel = "INFO"
)

// Default Logger configuration for DB.
const (
	defaultDBLogLevel = "INFO"
)
