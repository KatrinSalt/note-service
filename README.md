# Simple RESTful API for Note Taking with Azure Cosmos DB

## Description
Develop a simple RESTful API in Go that allows users to create, read, update, and delete notes. Use Azure Cosmos DB to store the notes. This project will introduce you to basic CRUD operations and integrating Azure Cosmos DB with a Go application.

## Main Functionality:
- Create new notes.
- Retrieve a list of notes.
- Update existing notes.
- Delete notes.

## Main Components:
1. **HTTP Server**: Set up using the Go `net/http` package.
2. **CRUD Handlers**: Implement handlers for creating, reading, updating, and deleting notes.
3. **Azure Cosmos DB Integration**: Use Azure Cosmos DB SDK for Go to store and retrieve notes.
4. **Configuration Management**: Manage connection strings and table storage settings.
5. **Error Handling**: Implement basic error handling for API requests.
