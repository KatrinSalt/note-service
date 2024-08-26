# Simple RESTful API for Note Taking with Azure Cosmos DB

## Description

This project is a simple RESTful API built in Go that allows users to create, read, update, and delete notes. The notes are stored in Azure Cosmos DB. This project demonstrates basic CRUD operations and the integration of Azure Cosmos DB with a Go application.

## Main Functionality

- **Create new notes**: Add a note under a specified category.
- **Retrieve a list of notes**: Retrieve all notes within a specified category.
- **Retrieve a note by ID**: Fetch a specific note by its ID within a category.
- **Update existing notes**: Modify the content of an existing note based on its ID and category.
- **Delete notes**: Remove a note based on its ID and category.

## API Endpoints

### Create a new note
- **Endpoint**: `POST /notes/create/{category}`
- **Description**: Creates a new note under the specified category.
- **Request Body**: 
    ```json
    {
        "note": "Note content here..."
    }
    ```

### Update an existing note
- **Endpoint**: `PUT /notes/update/{category}/{id}`
- **Description**: Updates an existing note identified by its ID and category.
- **Request Body**: 
    ```json
    {
        "note": "Updated note content here..."
    }
    ```

### Delete a note
- **Endpoint**: `DELETE /notes/delete/{category}/{id}`
- **Description**: Deletes a note identified by its ID and category.

### Retrieve a note by ID
- **Endpoint**: `GET /notes/categories/{category}/ids/{id}`
- **Description**: Retrieves a specific note by its ID within the specified category.

### Retrieve all notes in a category
- **Endpoint**: `GET /notes/categories/{category}`
- **Description**: Retrieves all notes within the specified category.

## CLI Client

In addition to the RESTful API, a CLI (Command Line Interface) client is available to interact with the API. The CLI allows users to create, read, update, and delete notes directly from the terminal.

### CLI Commands

- **Create a new note**:
    ```
    ./notes-service-cli create-note --category "category_name" --note "Note content here..."
    ```

- **Update an existing note**:
    ```
    ./notes-service-cli update-note --category "category_name" --id "note_id" --note "Updated note content here..."
    ```

- **Delete a note**:
    ```
    ./notes-service-cli delete-note --category "category_name" --id "note_id"
    ```

- **Retrieve a note by ID**:
    ```
    ./notes-service-cli get-note-by-id --category "category_name" --id "note_id"
    ```

- **Retrieve all notes in a category**:
    ```
    ./note-cli list-notes-by-category --category "category_name"
    ```

## Main Components

- **HTTP Server**: Set up using the Go `net/http` package.
- **CRUD Handlers**: Handlers for creating, reading, updating, and deleting notes.
- **Azure Cosmos DB Integration**: Utilizes Azure Cosmos DB SDK for Go to manage storage and retrieval of notes.
- **Error Handling**: Basic error handling for API requests.
- **Unit Tests**: Unit tests are included to ensure the correctness of the CRUD operations and other functionalities.

## Installation and Setup

1. Clone the repository:
    ```sh
    git clone https://github.com/KatrinSalt/note-service.git
    cd notes-service
    ```

2. Set up environment variables for Azure Cosmos DB:
    ```sh
    export COSMOSDB_CONNECTION_STRING="your-cosmos-db-connection-string"
    export COSMOSDB_DATABASE_ID="your-database-id"
    export COSMOSDB_CONTAINER_ID="your-container-id"
    export SERVICE_LOG_LEVEL="INFO"
    export DB_LOG_LEVEL="INFO"
    ```

3. Run the server:
    ```sh
    go run main.go
    ```

4. Use the CLI client to interact with the API:
    ```sh
    go build -o notes-service-cli cmd/cli/main.go
    ./notes-service-cli --help
    ```

## Conclusion

This project serves as a basic introduction to developing RESTful APIs with Go and integrating with Azure Cosmos DB. It is suitable for learning how to perform CRUD operations and managing cloud-based databases.
