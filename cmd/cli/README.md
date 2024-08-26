
# Notes Service CLI

`notes-service-cli` is a command-line interface (CLI) tool for managing notes through a simple CRUD (Create, Read, Update, Delete) interface with the notes service.

## Installation

To build and install the CLI tool, run the following commands:

```bash
go build -o notes-service-cli ./cmd/cli
```

This will create an executable file `notes-service-cli` in the current directory.

## Usage

The CLI supports several commands for creating, updating, deleting, and fetching notes from the notes service. Each command also supports various flags for customization.

### Global Flags

- `--host`, `-H`: The address of the service host. Default is `http://localhost:3000`.

### Commands

#### Create a Note

Creates a new note on the server.

**Usage:**

```bash
notes-service-cli create-note --category <category> --note <note content>
```

**Example:**

```bash
notes-service-cli create-note --category personal --note "Buy groceries"
notes-service-cli create -c work -n "Do time reporting"
```

#### Update a Note

Updates an existing note on the server.

**Usage:**

```bash
notes-service-cli update-note --category <category> --id <note id> --note <new note content>
```

**Example:**

```bash
notes-service-cli update-note --category personal --id 123 --note "Put groceries in the fridge"
notes-service-cli update -c work -i 321 -n "Do time reporting for the week 32"
```

#### Delete a Note

Deletes a note by ID on the server.

**Usage:**

```bash
notes-service-cli delete-note --category <category> --id <note id>
```

**Example:**

```bash
notes-service-cli delete-note --category personal --id 123
notes-service-cli delete -c work -i 321
```

#### Get a Note by ID

Fetches a note by category and ID from the server.

**Usage:**

```bash
notes-service-cli get-note-by-id --category <category> --id <note id>
```

**Example:**

```bash
notes-service-cli get-note-by-id --category personal --id 123
notes-service-cli get -c work -i 321
```

#### List Notes by Category

Lists all notes in a given category.

**Usage:**

```bash
notes-service-cli list-notes-by-category --category <category>
```

**Example:**

```bash
notes-service-cli list-notes-by-category --category personal
notes-service-cli list -c work
```

## Error Handling

The CLI provides error messages if something goes wrong during execution. This includes network errors, invalid inputs, or server errors. The errors are printed in red for easy identification.

## Contribution

Contributions are welcome! Please fork the repository and submit a pull request.
