package db

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"testing"

// 	"github.com/KatrinSalt/notes-service/log"

// 	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
// 	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// type mockInput struct {
// 	partitionKey azcosmos.PartitionKey
// 	item         []byte
// 	itemId       string
// 	o            *azcosmos.ItemOptions
// 	query        string
// }

// type mockCosmosContainerClient struct {
// 	t          *testing.T
// 	input      mockInput
// 	response   azcosmos.ItemResponse
// 	err        error
// 	funcCalled bool

// 	// pagerResponse *mockPager[azcosmos.QueryItemsResponse]
// 	pagerErr    error
// 	pagerCalled bool
// }

// // add response of type pager
// // or response type any[]

// func (m *mockCosmosContainerClient) CreateItem(ctx context.Context, partitionKey azcosmos.PartitionKey, item []byte, o *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
// 	m.funcCalled = true
// 	require.Equal(m.t, m.input.partitionKey, partitionKey)
// 	// require.Equal(m.t, m.input.item, item)
// 	require.Equal(m.t, m.input.o, o)
// 	return m.response, m.err
// }

// func (m *mockCosmosContainerClient) ReadItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, o *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
// 	return azcosmos.ItemResponse{}, nil
// }

// func (m *mockCosmosContainerClient) ReplaceItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, item []byte, o *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
// 	return azcosmos.ItemResponse{}, nil
// }

// func (m *mockCosmosContainerClient) DeleteItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, o *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
// 	return azcosmos.ItemResponse{}, nil
// }

// func (m *mockCosmosContainerClient) NewQueryItemsPager(query string, partitionKey azcosmos.PartitionKey, o *azcosmos.QueryOptions) *runtime.Pager[azcosmos.QueryItemsResponse] {
// 	require.Equal(m.t, m.input.query, query)
// 	require.Equal(m.t, m.input.partitionKey, partitionKey)
// 	require.Equal(m.t, m.input.o, o)
// 	return nil
// }

// type mockPager[T any] struct {
// 	Pages        []T   // Simulated pages of results
// 	CurrentPage  int   // Index of the current page
// 	MorePages    bool  // Determines if there are more pages
// 	ErrNextPage  error // Error to return from NextPage
// 	ErrUnmarshal error // Error to return from UnmarshalJSON
// }

// // More returns true if there are more pages to retrieve.
// func (m *mockPager[T]) More() bool {
// 	return m.MorePages && m.CurrentPage < len(m.Pages)-1
// }

// func (m *mockPager[T]) NextPage(ctx context.Context) (T, error) {
// 	if m.ErrNextPage != nil {
// 		// return error if set
// 		return *new(T), m.ErrNextPage
// 	}
// 	// simulate end of pages
// 	if m.CurrentPage >= len(m.Pages) {
// 		return *new(T), errors.New("no more pages")
// 	}

// 	result := m.Pages[m.CurrentPage]
// 	// advance to next page
// 	m.CurrentPage++
// 	m.MorePages = m.CurrentPage < len(m.Pages)

// 	return result, nil
// }

// func (m *mockPager[T]) UnmarshalJSON(data []byte) error {
// 	if m.ErrUnmarshal != nil {
// 		return m.ErrUnmarshal // return value if set
// 	}
// 	return json.Unmarshal(data, &m.Pages[m.CurrentPage])
// }

// func Test_CreateNote_Success(t *testing.T) {
// 	// Arrange
// 	mockClient := mockCosmosContainerClient{
// 		t: t,
// 		input: mockInput{
// 			partitionKey: azcosmos.NewPartitionKeyString("category"),
// 			item:         []byte("{\"id\":\"id\",\"category\":\"category\",\"note\":\"note\"}"),
// 			o: &azcosmos.ItemOptions{
// 				EnableContentResponseOnWrite: true,
// 			},
// 		},
// 		response: azcosmos.ItemResponse{
// 			Value: []byte(`{"test":"test"}`),
// 		},
// 		err: nil,
// 	}

// 	cosmosDB, err := NewCosmosDB(
// 		&mockClient,
// 		log.New(),
// 	)
// 	assert.NoError(t, err)

// 	note := Note{
// 		ID:       "id",
// 		Category: "category",
// 		Note:     "note",
// 	}

// 	// Act
// 	_, err = cosmosDB.CreateNote(context.Background(), note)

// 	// Assert
// 	require.NoError(t, err)
// 	require.True(t, mockClient.funcCalled)
// }

// func Test_CreateNote_Err_on_Create(t *testing.T) {
// 	// Arrange
// 	mockClient := mockCosmosContainerClient{
// 		t: t,
// 		input: mockInput{
// 			partitionKey: azcosmos.NewPartitionKeyString("category"),
// 			item:         []byte("{\"id\":\"id\",\"category\":\"category\",\"note\":\"note\"}"),
// 			o: &azcosmos.ItemOptions{
// 				EnableContentResponseOnWrite: true,
// 			},
// 		},
// 		response: azcosmos.ItemResponse{},
// 		err:      assert.AnError,
// 	}

// 	cosmosDB, err := NewCosmosDB(
// 		&mockClient,
// 		log.New(),
// 	)
// 	assert.NoError(t, err)

// 	note := Note{
// 		ID:       "id",
// 		Category: "category",
// 		Note:     "note",
// 	}

// 	// Act
// 	_, err = cosmosDB.CreateNote(context.Background(), note)

// 	// Assert
// 	require.Error(t, err)
// 	require.True(t, mockClient.funcCalled)
// }

// func Test_CreateNote_Err_on_NotJSONResponse(t *testing.T) {
// 	// Arrange
// 	mockClient := mockCosmosContainerClient{
// 		t: t,
// 		input: mockInput{
// 			partitionKey: azcosmos.NewPartitionKeyString("category"),
// 			item:         []byte("{\"id\":\"id\",\"category\":\"category\",\"note\":\"note\"}"),
// 			o: &azcosmos.ItemOptions{
// 				EnableContentResponseOnWrite: true,
// 			},
// 		},
// 		response: azcosmos.ItemResponse{
// 			Value: []byte(`notajson`),
// 		},
// 		err: nil,
// 	}

// 	cosmosDB, err := NewCosmosDB(
// 		&mockClient,
// 		log.New(),
// 	)
// 	assert.NoError(t, err)

// 	note := Note{
// 		ID:       "id",
// 		Category: "category",
// 		Note:     "note",
// 	}

// 	// Act
// 	_, err = cosmosDB.CreateNote(context.Background(), note)

// 	// Assert
// 	require.Error(t, err)
// 	require.True(t, mockClient.funcCalled)
// }

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
