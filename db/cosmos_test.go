package db

import (
	"context"
	"testing"
	"time"

	"github.com/KatrinSalt/notes-service/log"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockInput struct {
	ctx          context.Context
	partitionKey azcosmos.PartitionKey
	item         []byte
	itemId       string
	query        string
	o            *azcosmos.ItemOptions
}

type mockCosmosContainerClient struct {
	t          *testing.T
	input      mockInput
	response   []byte
	err        error
	funcCalled bool
}

func (m *mockCosmosContainerClient) CreateItem(ctx context.Context, partitionKey azcosmos.PartitionKey, item []byte, o *azcosmos.ItemOptions) ([]byte, error) {
	m.funcCalled = true

	require.Equal(m.t, m.input.ctx, ctx)
	require.Equal(m.t, m.input.partitionKey, partitionKey)
	require.Equal(m.t, m.input.item, item)
	require.Equal(m.t, m.input.o, o)

	return m.response, m.err
}

func (m *mockCosmosContainerClient) ReadItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, o *azcosmos.ItemOptions) ([]byte, error) {
	return nil, nil
}

func (m *mockCosmosContainerClient) ReplaceItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, item []byte, o *azcosmos.ItemOptions) ([]byte, error) {
	return nil, nil
}

func (m *mockCosmosContainerClient) DeleteItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, o *azcosmos.ItemOptions) ([]byte, error) {
	return nil, nil
}

func (m *mockCosmosContainerClient) Query(ctx context.Context, query string, partitionKey azcosmos.PartitionKey, o *azcosmos.QueryOptions) ([][]byte, error) {
	require.Equal(m.t, m.input.query, query)
	require.Equal(m.t, m.input.partitionKey, partitionKey)
	require.Equal(m.t, m.input.o, o)
	return nil, nil
}

func Test_CreateNote_Success(t *testing.T) {
	mockID := "123e4567-e89b-12d3-a456-426614174000"
	mockCreatedAt := time.Now().UTC()

	// Arrange
	mockClient := mockCosmosContainerClient{
		t: t,
		input: mockInput{
			ctx:          context.Background(),
			partitionKey: azcosmos.NewPartitionKeyString("category"),
			item:         []byte("{\"id\":\"" + mockID + "\",\"category\":\"category\",\"note\":\"note\",\"timestamp\":\"" + mockCreatedAt.Format(time.RFC3339Nano) + "\"}"),
			o: &azcosmos.ItemOptions{
				EnableContentResponseOnWrite: true,
			},
		},
		response: []byte("{\"id\":\"" + mockID + "\",\"category\":\"category\",\"note\":\"note\",\"timestamp\":\"" + mockCreatedAt.Format(time.RFC3339Nano) + "\"}"),
		err:      nil,
	}

	cosmosDB, err := NewNotesDB(
		&mockClient,
		log.New(),
	)
	assert.NoError(t, err)

	inputNote := Note{
		ID:        mockID,
		Category:  "category",
		Note:      "note",
		CreatedAt: mockCreatedAt,
	}

	// Act
	noteDB, err := cosmosDB.CreateNote(context.Background(), inputNote)

	// Assert
	require.NoError(t, err)
	require.True(t, mockClient.funcCalled)
	require.NotEmpty(t, noteDB.ID)
	require.Equal(t, inputNote.Category, noteDB.Category)
	require.Equal(t, inputNote.Note, noteDB.Note)
	// Check that CreatedAt is within a reasonable range
	expectedTime := time.Now().UTC()
	require.WithinDuration(t, expectedTime, noteDB.CreatedAt, time.Second*2)
}

func Test_CreateNote_Err_on_Create(t *testing.T) {
	mockID := "123e4567-e89b-12d3-a456-426614174000"
	mockCreatedAt := time.Now().UTC()

	// Arrange
	mockClient := mockCosmosContainerClient{
		t: t,
		input: mockInput{
			ctx:          context.Background(),
			partitionKey: azcosmos.NewPartitionKeyString("category"),
			item:         []byte("{\"id\":\"" + mockID + "\",\"category\":\"category\",\"note\":\"note\",\"timestamp\":\"" + mockCreatedAt.Format(time.RFC3339Nano) + "\"}"),
			o: &azcosmos.ItemOptions{
				EnableContentResponseOnWrite: true,
			},
		},
		response: nil,
		err:      assert.AnError,
	}

	cosmosDB, err := NewNotesDB(
		&mockClient,
		log.New(),
	)
	assert.NoError(t, err)

	inputNote := Note{
		ID:        mockID,
		Category:  "category",
		Note:      "note",
		CreatedAt: mockCreatedAt,
	}

	// Act
	noteDB, err := cosmosDB.CreateNote(context.Background(), inputNote)

	// Assert
	require.True(t, mockClient.funcCalled)
	require.Equal(t, Note{}, noteDB)
	require.Error(t, err)
}

func Test_CreateNote_Err_on_NotJSONResponse(t *testing.T) {
	mockID := "123e4567-e89b-12d3-a456-426614174000"
	mockCreatedAt := time.Now().UTC()

	// Arrange
	mockClient := mockCosmosContainerClient{
		t: t,
		input: mockInput{
			ctx:          context.Background(),
			partitionKey: azcosmos.NewPartitionKeyString("category"),
			item:         []byte("{\"id\":\"" + mockID + "\",\"category\":\"category\",\"note\":\"note\",\"timestamp\":\"" + mockCreatedAt.Format(time.RFC3339Nano) + "\"}"),
			o: &azcosmos.ItemOptions{
				EnableContentResponseOnWrite: true,
			},
		},
		response: []byte(`notajson`),
		err:      nil,
	}

	cosmosDB, err := NewNotesDB(
		&mockClient,
		log.New(),
	)
	assert.NoError(t, err)

	inputNote := Note{
		ID:        mockID,
		Category:  "category",
		Note:      "note",
		CreatedAt: mockCreatedAt,
	}

	// Act
	noteDB, err := cosmosDB.CreateNote(context.Background(), inputNote)

	// Assert
	require.True(t, mockClient.funcCalled)
	require.Equal(t, Note{}, noteDB)
	require.Error(t, err)
}

// func Test_GetNotesByCategory_Success(t *testing.T) {
// 	mockClient := mockCosmosContainerClient{
// 		t: t,
// 		input: mockInput{
// 			partitionKey: azcosmos.NewPartitionKeyString("category"),
// 			query:        "SELECT * FROM c",
// 			o:            nil,
// 		},
// 		err: nil,
// 	}

// 	cosmosDB, err := NewCosmosDB(
// 		&mockClient,
// 		log.New(),
// 	)
// 	assert.NoError(t, err)

// 	var notes []Note

// 	notes, err = cosmosDB.GetNotesByCategory(context.Background(), "category")

// 	require.Error(t, err)

// 	require.Greater(t, len(notes), 0)

// 	notesExpexted := []Note{
// 		{
// 			ID:       "1",
// 			Category: "category",
// 			Note:     "first note",
// 		},
// 		{
// 			ID:       "2",
// 			Category: "category",
// 			Note:     "second note",
// 		},
// 	}

// 	require.Equal(t, notesExpexted, notes)

// }
