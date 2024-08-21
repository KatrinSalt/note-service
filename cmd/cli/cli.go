package main

import (
	"github.com/KatrinSalt/notes-service/cmd/cli/commands"
	"github.com/KatrinSalt/notes-service/cmd/cli/output"
	"github.com/urfave/cli/v2"
)

const (
	name = "notes-service-cli"
)

var host string

func CLI(args []string) int {
	app := &cli.App{
		Name:  name,
		Usage: "A CLI to handle CRUD operations with the note service.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Aliases:     []string{"H"},
				Usage:       "The address of the service host",
				Value:       "http://localhost:3000", // Default value
				Destination: &host,
			},
		},
		Commands: []*cli.Command{
			commands.CreateNote(&host),
			commands.UpdateNote(&host),
			commands.DeleteNote(&host),
			commands.GetNoteByID(&host),
			commands.ListNotes(&host),
		},
		CustomAppHelpTemplate: `NAME:
	{{.HelpName}} - {{.Usage}}

USAGE:
   	{{.HelpName}} [global options] command [command options]

GLOBAL OPTIONS:
   	{{range .VisibleFlags}}{{.}}
   	{{end}}

COMMANDS:
   	{{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   	{{end}}

{{if .Commands}}
COMMAND DETAILS:
{{range .Commands}}
	{{.Name}}, {{join .Aliases ", "}}

    OPTIONS:
	  {{range .VisibleFlags}}    {{.}}
	  {{end}}{{if .UsageText}}
		EXAMPLE:
	  {{.UsageText}}{{end}}

{{end}}{{end}}`,
	}

	if err := app.Run(args); err != nil {
		output.PrintlnErr(err)
		return 1
	}
	return 0
}
