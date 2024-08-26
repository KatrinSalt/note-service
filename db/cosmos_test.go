package db

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateNote(t *testing.T) {
	mockID := "123e4567-e89b-12d3-a456-426614174000"
	mockCreatedAt := time.Now().UTC()

	tests := []struct {
		name          string
		inputNote     Note
		mockResponse  []byte
		mockError     error
		expectedNote  Note
		expectError   bool
		expectedError error
	}{
		{
			name: "CreateNote() - successful creation",
			inputNote: Note{
				ID:        mockID,
				Category:  "category",
				Note:      "note",
				CreatedAt: mockCreatedAt,
			},
			mockResponse: []byte("{\"id\":\"" + mockID + "\",\"category\":\"category\",\"note\":\"note\",\"timestamp\":\"" + mockCreatedAt.Format(time.RFC3339Nano) + "\"}"),
			mockError:    nil,
			expectedNote: Note{
				ID:        mockID,
				Category:  "category",
				Note:      "note",
				CreatedAt: mockCreatedAt,
			},
			expectError:   false,
			expectedError: nil,
		},
		{
			name: "CreateNote() - internal db error",
			inputNote: Note{
				ID:        mockID,
				Category:  "category",
				Note:      "note",
				CreatedAt: mockCreatedAt,
			},
			mockResponse:  nil,
			mockError:     assert.AnError,
			expectedNote:  Note{},
			expectError:   true,
			expectedError: fmt.Errorf("%w: %w", ErrInternalDB, assert.AnError),
		},
		{
			name: "CreateNote() - error on non-json response",
			inputNote: Note{
				ID:        mockID,
				Category:  "category",
				Note:      "note",
				CreatedAt: mockCreatedAt,
			},
			mockResponse:  []byte(`notajson`),
			mockError:     nil,
			expectedNote:  Note{},
			expectError:   true,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockClient := mockCosmosContainerClient{
				t: t,
				input: mockInput{
					ctx:          context.Background(),
					partitionKey: tt.inputNote.Category,
					item:         []byte("{\"id\":\"" + tt.inputNote.ID + "\",\"category\":\"" + tt.inputNote.Category + "\",\"note\":\"" + tt.inputNote.Note + "\",\"timestamp\":\"" + mockCreatedAt.Format(time.RFC3339Nano) + "\"}"),
					id:           tt.inputNote.ID,
				},
				response: tt.mockResponse,
				err:      tt.mockError,
			}

			cosmosDB, err := NewNotesDB(
				&mockClient,
			)
			assert.NoError(t, err)

			// Act
			noteDB, err := cosmosDB.CreateNote(context.Background(), tt.inputNote)

			// Assert
			require.True(t, mockClient.funcCalled)

			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, tt.expectedNote, noteDB)
				if tt.expectedError != nil {
					require.Equal(t, tt.expectedError, err)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedNote.ID, noteDB.ID)
				require.Equal(t, tt.expectedNote.Category, noteDB.Category)
				require.Equal(t, tt.expectedNote.Note, noteDB.Note)
				require.WithinDuration(t, tt.expectedNote.CreatedAt, noteDB.CreatedAt, time.Second*2)
			}
		})
	}
}

func Test_UpdateNote(t *testing.T) {
	mockID := "123e4567-e89b-12d3-a456-426614174000"
	mockCreatedAt := time.Now().UTC()

	tests := []struct {
		name          string
		inputNote     Note
		mockResponse  []byte
		mockError     error
		expectedNote  Note
		expectError   bool
		expectedError error
	}{
		{
			name: "UpdateNote() - successful update",
			inputNote: Note{
				ID:        mockID,
				Category:  "category",
				Note:      "updated note",
				CreatedAt: mockCreatedAt,
			},
			mockResponse: []byte("{\"id\":\"" + mockID + "\",\"category\":\"category\",\"note\":\"updated note\",\"timestamp\":\"" + mockCreatedAt.Format(time.RFC3339Nano) + "\"}"),
			mockError:    nil,
			expectedNote: Note{
				ID:        mockID,
				Category:  "category",
				Note:      "updated note",
				CreatedAt: mockCreatedAt,
			},
			expectError:   false,
			expectedError: nil,
		},
		{
			name: "UpdateNote() - internal db error",
			inputNote: Note{
				ID:        mockID,
				Category:  "category",
				Note:      "updated note",
				CreatedAt: mockCreatedAt,
			},
			mockResponse:  nil,
			mockError:     assert.AnError,
			expectedNote:  Note{},
			expectError:   true,
			expectedError: fmt.Errorf("%w: %w", ErrInternalDB, assert.AnError),
		},
		{
			name: "UpdateNote() - not found error",
			inputNote: Note{
				ID:        mockID,
				Category:  "category",
				Note:      "updated note",
				CreatedAt: mockCreatedAt,
			},
			mockResponse: nil,
			mockError: &azcore.ResponseError{
				ErrorCode:   "Resource not found",
				StatusCode:  http.StatusNotFound,
				RawResponse: nil,
			},
			expectedNote:  Note{},
			expectError:   true,
			expectedError: ErrNotFound,
		},
		{
			name: "UpdateNote() - error on non-json response",
			inputNote: Note{
				ID:        mockID,
				Category:  "category",
				Note:      "updated note",
				CreatedAt: mockCreatedAt,
			},
			mockResponse:  []byte(`notajson`),
			mockError:     nil,
			expectedNote:  Note{},
			expectError:   true,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockClient := mockCosmosContainerClient{
				t: t,
				input: mockInput{
					ctx:          context.Background(),
					partitionKey: tt.inputNote.Category,
					item:         []byte("{\"id\":\"" + tt.inputNote.ID + "\",\"category\":\"" + tt.inputNote.Category + "\",\"note\":\"" + tt.inputNote.Note + "\",\"timestamp\":\"" + mockCreatedAt.Format(time.RFC3339Nano) + "\"}"),
					id:           tt.inputNote.ID,
				},
				response: tt.mockResponse,
				err:      tt.mockError,
			}

			cosmosDB, err := NewNotesDB(
				&mockClient,
			)
			assert.NoError(t, err)

			// Act
			noteDB, err := cosmosDB.UpdateNote(context.Background(), tt.inputNote)

			// Assert
			require.True(t, mockClient.funcCalled)

			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, tt.expectedNote, noteDB)
				if tt.expectedError != nil {
					require.Equal(t, tt.expectedError, err)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedNote.ID, noteDB.ID)
				require.Equal(t, tt.expectedNote.Category, noteDB.Category)
				require.Equal(t, tt.expectedNote.Note, noteDB.Note)
				require.WithinDuration(t, tt.expectedNote.CreatedAt, noteDB.CreatedAt, time.Second*2)
			}
		})
	}
}

func Test_DeleteNote(t *testing.T) {
	mockID := "123e4567-e89b-12d3-a456-426614174000"
	mockCreatedAt := time.Now().UTC()

	tests := []struct {
		name          string
		inputNote     Note
		mockResponse  []byte
		mockError     error
		expectedNote  Note
		expectError   bool
		expectedError error
	}{
		{
			name: "DeleteNote() - successful deletion",
			inputNote: Note{
				ID:        mockID,
				Category:  "category",
				Note:      "note",
				CreatedAt: mockCreatedAt,
			},
			mockError:     nil,
			expectError:   false,
			expectedError: nil,
		},
		{
			name: "DeleteNote() - internal db error",
			inputNote: Note{
				ID:        mockID,
				Category:  "category",
				Note:      "note",
				CreatedAt: mockCreatedAt,
			},
			mockError:     assert.AnError,
			expectError:   true,
			expectedError: fmt.Errorf("%w: %w", ErrInternalDB, assert.AnError),
		},
		{
			name: "DeleteNote() - not found error",
			inputNote: Note{
				ID:        mockID,
				Category:  "category",
				Note:      "note",
				CreatedAt: mockCreatedAt,
			},
			mockError: &azcore.ResponseError{
				ErrorCode:   "Resource not found",
				StatusCode:  http.StatusNotFound,
				RawResponse: nil,
			},
			expectError:   true,
			expectedError: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockClient := mockCosmosContainerClient{
				t: t,
				input: mockInput{
					ctx:          context.Background(),
					partitionKey: tt.inputNote.Category,
					item:         []byte("{\"id\":\"" + tt.inputNote.ID + "\",\"category\":\"" + tt.inputNote.Category + "\",\"note\":\"" + tt.inputNote.Note + "\",\"timestamp\":\"" + mockCreatedAt.Format(time.RFC3339Nano) + "\"}"),
					id:           tt.inputNote.ID,
				},
				response: tt.mockResponse,
				err:      tt.mockError,
			}

			cosmosDB, err := NewNotesDB(
				&mockClient,
			)
			assert.NoError(t, err)

			// Act
			err = cosmosDB.DeleteNote(context.Background(), tt.inputNote.ID, tt.inputNote.Category)

			// Assert
			require.True(t, mockClient.funcCalled)

			if tt.expectError {
				require.Error(t, err)
				if tt.expectedError != nil {
					require.Equal(t, tt.expectedError, err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_GetNotesByCategory(t *testing.T) {
	mockCategory := "test"
	mockID1 := "123e4567-e89b-12d3-a456-426614174000"
	mockCreatedAt1 := time.Now().UTC()

	mockID2 := "f4d2e990-8f5d-4c63-85ff-c8e2b4d23b8a"
	mockCreatedAt2 := time.Now().UTC()

	mockByteSlice1 := []byte("{\"id\":\"" + mockID1 + "\",\"category\":\"" + mockCategory + "\",\"note\":\"test note 1\",\"timestamp\":\"" + mockCreatedAt1.Format(time.RFC3339Nano) + "\"}")
	mockByteSlice2 := []byte("{\"id\":\"" + mockID2 + "\",\"category\":\"" + mockCategory + "\",\"note\":\"test note 2\",\"timestamp\":\"" + mockCreatedAt2.Format(time.RFC3339Nano) + "\"}")

	mockResponseByteSlices := [][]byte{mockByteSlice1, mockByteSlice2}

	tests := []struct {
		name          string
		inputCategory string
		mockResponse  [][]byte
		mockError     error
		expectedNotes []Note
		expectError   bool
		expectedError error
	}{
		{
			name:          "GetNotesByCategory() - successful execution",
			inputCategory: mockCategory,
			mockResponse:  mockResponseByteSlices,
			mockError:     nil,
			expectedNotes: []Note{
				{
					ID:        mockID1,
					Category:  mockCategory,
					Note:      "test note 1",
					CreatedAt: mockCreatedAt1,
				},
				{
					ID:        mockID2,
					Category:  mockCategory,
					Note:      "test note 2",
					CreatedAt: mockCreatedAt2,
				},
			},
			expectError:   false,
			expectedError: nil,
		},
		{
			name:          "GetNotesByCategory() - no items found in provided category",
			inputCategory: mockCategory,
			mockResponse:  [][]byte{},
			mockError:     nil,
			expectedNotes: []Note(nil),
			expectError:   false,
			expectedError: nil,
		},
		{
			name:          "GetNotesByCategory() - internal db error",
			inputCategory: mockCategory,
			mockResponse:  nil,
			mockError:     assert.AnError,
			expectedNotes: []Note{},
			expectError:   true,
			expectedError: fmt.Errorf("%w: %w", ErrInternalDB, assert.AnError),
		},
		{
			name:          "GetNotesByCategory() - error on non-json response",
			inputCategory: mockCategory,
			mockResponse:  [][]byte{[]byte(`notajson`), []byte(`notajson2`)},
			mockError:     nil,
			expectedNotes: []Note{},
			expectError:   true,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockClient := mockCosmosContainerClient{
				t: t,
				input: mockInput{
					ctx:          context.Background(),
					partitionKey: tt.inputCategory,
				},
				responses: tt.mockResponse,
				err:       tt.mockError,
			}

			cosmosDB, err := NewNotesDB(
				&mockClient,
			)
			assert.NoError(t, err)

			// Act
			notesDB, err := cosmosDB.GetNotesByCategory(context.Background(), tt.inputCategory)

			// Assert
			require.True(t, mockClient.funcCalled)

			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, tt.expectedNotes, notesDB)
				if tt.expectedError != nil {
					require.Equal(t, tt.expectedError, err)
				}
			} else {
				require.NoError(t, err)
				require.Len(t, notesDB, len(tt.expectedNotes))
				if len(notesDB) > 0 {
					for i, noteDB := range notesDB {
						require.Equal(t, tt.expectedNotes[i].ID, noteDB.ID)
						require.Equal(t, tt.expectedNotes[i].Category, noteDB.Category)
						require.Equal(t, tt.expectedNotes[i].Note, noteDB.Note)
						require.WithinDuration(t, tt.expectedNotes[i].CreatedAt, noteDB.CreatedAt, time.Second*2)
					}
				} else {
					require.Equal(t, tt.expectedNotes, notesDB)
				}
			}
		})
	}
}

func Test_GetNoteByID(t *testing.T) {
	mockCategory := "test"
	mockID := "123e4567-e89b-12d3-a456-426614174000"
	mockCreatedAt := time.Now().UTC()

	tests := []struct {
		name          string
		inputID       string
		inputCategory string
		mockResponse  []byte
		mockError     error
		expectedNote  Note
		expectError   bool
		expectedError error
	}{
		{
			name:          "GetNoteByID() - successful execution",
			inputID:       mockID,
			inputCategory: mockCategory,
			mockResponse:  []byte("{\"id\":\"" + mockID + "\",\"category\":\"" + mockCategory + "\",\"note\":\"test note\",\"timestamp\":\"" + mockCreatedAt.Format(time.RFC3339Nano) + "\"}"),
			mockError:     nil,
			expectedNote: Note{
				ID:        "123e4567-e89b-12d3-a456-426614174000",
				Category:  "test",
				Note:      "test note",
				CreatedAt: mockCreatedAt,
			},
			expectError:   false,
			expectedError: nil,
		},
		{
			name:          "GetNoteByID() - internal db error",
			inputID:       mockID,
			inputCategory: mockCategory,
			mockResponse:  nil,
			mockError:     assert.AnError,
			expectedNote:  Note{},
			expectError:   true,
			expectedError: fmt.Errorf("%w: %w", ErrInternalDB, assert.AnError),
		},
		{
			name:          "GetNoteByID() - not found error",
			inputID:       mockID,
			inputCategory: mockCategory,
			mockResponse:  nil,
			mockError: &azcore.ResponseError{
				ErrorCode:   "Resource not found",
				StatusCode:  http.StatusNotFound,
				RawResponse: nil,
			},
			expectedNote:  Note{},
			expectError:   true,
			expectedError: ErrNotFound,
		},

		{
			name:          "GetNoteByID() - error on non-json response",
			inputID:       mockID,
			inputCategory: mockCategory,
			mockResponse:  []byte(`notajson`),
			mockError:     nil,
			expectedNote:  Note{},
			expectError:   true,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockClient := mockCosmosContainerClient{
				t: t,
				input: mockInput{
					ctx:          context.Background(),
					partitionKey: tt.inputCategory,
					id:           tt.inputID,
				},
				response: tt.mockResponse,
				err:      tt.mockError,
			}

			cosmosDB, err := NewNotesDB(
				&mockClient,
			)
			assert.NoError(t, err)

			// Act
			noteDB, err := cosmosDB.GetNoteByID(context.Background(), tt.inputCategory, tt.inputID)

			// Assert
			require.True(t, mockClient.funcCalled)

			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, tt.expectedNote, noteDB)
				if tt.expectedError != nil {
					require.Equal(t, tt.expectedError, err)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedNote.ID, noteDB.ID)
				require.Equal(t, tt.expectedNote.Category, noteDB.Category)
				require.Equal(t, tt.expectedNote.Note, noteDB.Note)
				require.WithinDuration(t, tt.expectedNote.CreatedAt, noteDB.CreatedAt, time.Second*2)
			}
		})
	}
}

type mockInput struct {
	ctx          context.Context
	partitionKey string
	item         []byte
	id           string
}

type mockCosmosContainerClient struct {
	t          *testing.T
	input      mockInput
	response   []byte
	responses  [][]byte
	err        error
	funcCalled bool
}

func (m *mockCosmosContainerClient) CreateItem(ctx context.Context, partitionKey string, item []byte) ([]byte, error) {
	m.funcCalled = true

	require.Equal(m.t, m.input.ctx, ctx)
	require.Equal(m.t, m.input.partitionKey, partitionKey)
	require.Equal(m.t, m.input.item, item)

	return m.response, m.err
}

func (m *mockCosmosContainerClient) ReadItem(ctx context.Context, partitionKey string, id string) ([]byte, error) {
	m.funcCalled = true

	require.Equal(m.t, m.input.ctx, ctx)
	require.Equal(m.t, m.input.partitionKey, partitionKey)
	require.Equal(m.t, m.input.id, id)

	return m.response, m.err
}

func (m *mockCosmosContainerClient) ReplaceItem(ctx context.Context, partitionKey string, id string, item []byte) ([]byte, error) {
	m.funcCalled = true

	require.Equal(m.t, m.input.ctx, ctx)
	require.Equal(m.t, m.input.partitionKey, partitionKey)
	require.Equal(m.t, m.input.id, id)
	require.Equal(m.t, m.input.item, item)

	return m.response, m.err
}

func (m *mockCosmosContainerClient) DeleteItem(ctx context.Context, partitionKey string, id string) error {
	m.funcCalled = true

	require.Equal(m.t, m.input.ctx, ctx)
	require.Equal(m.t, m.input.partitionKey, partitionKey)
	require.Equal(m.t, m.input.id, id)

	return m.err
}

func (m *mockCosmosContainerClient) ListItems(ctx context.Context, partitionKey string) ([][]byte, error) {
	m.funcCalled = true

	require.Equal(m.t, m.input.partitionKey, partitionKey)

	return m.responses, m.err
}
