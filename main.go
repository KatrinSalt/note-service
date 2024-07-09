package main

import (
	"fmt"
	"os"

	"github.com/KatrinSalt/notes-service/config"
	"github.com/KatrinSalt/notes-service/server"
)

func main() {
	fmt.Println("Starting the note service...")
	cfg, err := config.New()
	if err != nil {
		fmt.Printf("could not load configuration: %s\n", err)
		os.Exit(1)
	}

	// DEBUG: Print the configuration struct
	fmt.Printf("Configuration: %+v\n", cfg)

	services, err := config.SetupServices(cfg.Services)
	if err != nil {
		fmt.Printf("could not setup services: %s\n", err)
		os.Exit(1)
	}

	srv, err := server.New(
		services.Note,
		server.WithAddress(cfg.Server.Host+":"+cfg.Server.Port),
	)
	if err != nil {
		fmt.Printf("could not create server: %s\n", err)
		os.Exit(1)
	}

	if err := srv.Start(); err != nil {
		fmt.Printf("could not start server: %s\n", err)
		os.Exit(1)
	}
}
