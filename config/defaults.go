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
	defaultCosmosDatabaseID  = "NotesDB"
	defaultCosmosContainerID = "notes"
)
