package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/KatrinSalt/notes-service/cmd/cli/output"
	"github.com/urfave/cli/v2"
)

type Note struct {
	ID       string `json:"id,omitempty"`
	Category string `json:"category,omitempty"`
	Note     string `json:"note,omitempty"`
}

type Response struct {
	Message string `json:"message,omitempty"`
	Note    Note   `json:"note,omitempty"`
	Notes   []Note `json:"notes,omitempty"`
}

func CreateNote(host *string) *cli.Command {
	return &cli.Command{
		Name:    "create-note",
		Aliases: []string{"create"},
		Usage:   "Create a new note on the server",
		UsageText: ` 
		    notes-service-cli create-note --category personal --note "Buy groceries"
		    notes-service-cli create -c work -n "Do time reporting"`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "category",
				Aliases:  []string{"c"},
				Usage:    "Category of the note, required",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "note",
				Aliases:  []string{"n"},
				Usage:    "Content of the note to create",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			category := c.String("category")
			noteContent := c.String("note")

			jsonStr := []byte(fmt.Sprintf(`{"note":"%s"}`, noteContent))
			url := fmt.Sprintf("%s/notes/create/%s", *host, category)
			reqResp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
			if err != nil {
				return fmt.Errorf("error creating note: %w", err)
			}
			defer reqResp.Body.Close()

			response, err := processResponse(reqResp)
			if err != nil {
				return fmt.Errorf("error creating the note: %w", err)
			}

			note := Note{}
			if response.Note == note {
				output.Println(response.Message)
			} else {
				message := fmt.Sprintf("Note is created.\nNote Details:\n  ID: %s\n  Category: %s\n  Note: %s", response.Note.ID, response.Note.Category, response.Note.Note)
				output.Println(message)
			}
			return nil
		},
	}
}

func UpdateNote(host *string) *cli.Command {
	return &cli.Command{
		Name:    "update-note",
		Aliases: []string{"update"},
		Usage:   "Update an existing note on the server",
		UsageText: ` 
        notes-service-cli update-note --category personal --id 123 --note "Put groceries in the fridge"
        notes-service-cli update -c work -i 321 -n "Do time reporting for the week 32"`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "category",
				Aliases:  []string{"c"},
				Usage:    "Category of the note, required",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "ID of the note to update, required",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "note",
				Aliases:  []string{"n"},
				Usage:    "New content of the note",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			category := c.String("category")
			id := c.String("id")
			noteContent := c.String("note")

			jsonStr := []byte(fmt.Sprintf(`{"note":"%s"}`, noteContent))

			url := fmt.Sprintf("%s/notes/update/%s/%s", *host, category, id)
			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonStr))
			if err != nil {
				return fmt.Errorf("error creating update request: %w", err)
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			reqResp, err := client.Do(req)
			if err != nil {
				// fmt.Println("Error updating note:", err)
				return fmt.Errorf("error updating note: %w", err)
			}
			defer reqResp.Body.Close()

			response, err := processResponse(reqResp)
			if err != nil {
				return fmt.Errorf("error updating the note: %w", err)
			}

			note := Note{}
			if response.Note == note {
				output.Println(response.Message)
			} else {
				message := fmt.Sprintf("Note is updated.\nNote Details:\n  ID: %s\n  Category: %s\n  Note: %s", response.Note.ID, response.Note.Category, response.Note.Note)
				output.Println(message)
			}
			return nil
		},
	}
}

func DeleteNote(host *string) *cli.Command {
	return &cli.Command{
		Name:    "delete-note",
		Aliases: []string{"delete"},
		Usage:   "Delete a note by ID on the server",
		UsageText: ` 
        notes-service-cli delete-note --category personal --id 123
        notes-service-cli delete -c work -i 321`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "category",
				Aliases:  []string{"c"},
				Usage:    "Category of the note to delete, required",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "ID of the note to delete, required",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			category := c.String("category")
			id := c.String("id")

			url := fmt.Sprintf("%s/notes/delete/%s/%s", *host, category, id)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			if err != nil {
				return fmt.Errorf("error creating delete request: %w", err)
			}

			client := &http.Client{}
			reqResp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("error deleting the note: %w", err)
			}
			defer reqResp.Body.Close()

			response, err := processResponse(reqResp)
			if err != nil {
				return fmt.Errorf("error deleting note: %w", err)
			}

			note := Note{}
			if response.Note == note {
				output.Println(response.Message)
			} else {
				message := fmt.Sprintf("Note is deleted.\nNote Details:\n  ID: %s\n  Category: %s\n  Note: %s", response.Note.ID, response.Note.Category, response.Note.Note)
				output.Println(message)
			}
			return nil
		},
	}
}

func GetNoteByID(host *string) *cli.Command {
	return &cli.Command{
		Name:    "get-note-by-id",
		Aliases: []string{"get"},
		Usage:   "Fetch a note by category and ID from the server",
		UsageText: ` 
        notes-service-cli get-note-by-id --category personal --id 123
        notes-service-cli get -c work -i 321`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "category",
				Aliases:  []string{"c"},
				Usage:    "Category of the note to fetch, required",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "ID of the note to fetch, required",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			category := c.String("category")
			id := c.String("id")

			if category == "" {
				// fmt.Println("Please provide category of the note to delete.")
				return fmt.Errorf("note category shall be provided")
			}

			if id == "" {
				// fmt.Println("Please provide ID of the note to delete.")
				return fmt.Errorf("note ID shall be provided")
			}

			url := fmt.Sprintf("%s/notes/categories/%s/ids/%s", *host, category, id)

			reqResp, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("error fetching the note: %w", err)
			}
			defer reqResp.Body.Close()

			response, err := processResponse(reqResp)
			if err != nil {
				return fmt.Errorf("error fetching the note: %w", err)
			}

			note := Note{}
			if response.Note == note {
				output.Println(response.Message)
			} else {
				message := fmt.Sprintf("Note is fetched.\nNote Details:\n  ID: %s\n  Category: %s\n  Note: %s", response.Note.ID, response.Note.Category, response.Note.Note)
				output.Println(message)
			}
			return nil
		},
	}
}

func ListNotes(host *string) *cli.Command {
	return &cli.Command{
		Name:    "list-notes-by-category",
		Aliases: []string{"list"},
		Usage:   "List notes by category from the server",
		UsageText: ` 
        notes-service-cli list-notes-by-category --category personal
        notes-service-cli list -c work`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "category",
				Aliases:  []string{"c"},
				Usage:    "Category of the notes to list, required",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			category := c.String("category")

			if category == "" {
				// fmt.Println("Please provide category of the note to delete.")
				return fmt.Errorf("note category shall be provided")
			}

			url := fmt.Sprintf("%s/notes/categories/%s", *host, category)

			reqResp, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("error listing the notes: %w", err)
			}
			defer reqResp.Body.Close()

			response, err := processResponse(reqResp)
			if err != nil {
				return fmt.Errorf("error listing the notes: %w", err)
			}

			if len(response.Notes) == 0 {
				message := fmt.Sprintf("No notes found in the category '%s'.", category)
				output.Println(message)
			} else {
				message := fmt.Sprintf("List of the notes in the category '%s':", category)
				output.Println(message)
				for _, note := range response.Notes {
					noteStr := fmt.Sprintf("ID: %s | Note: %s", note.ID, note.Note)
					output.Println(noteStr)
				}
			}
			return nil
		},
	}
}

// processResponse checks the HTTP status code, reads and unmarshals the response body,
// and returns a formatted message depending on the operation and outcome.
func processResponse(resp *http.Response) (Response, error) {
	// Check if the status code indicates an error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return Response{}, fmt.Errorf("status: %s, response: %s", resp.Status, string(bodyBytes))
	}

	// Read and process the response body for successful requests
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return Response{}, err
	}

	return response, nil
}
