package db

import (
	"context"
	"testing"

	"github.com/KatrinSalt/notes-service/log"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockInput struct {
	partitionKey azcosmos.PartitionKey
	item         []byte
	itemId       string
	o            *azcosmos.ItemOptions
}

type mockCosmosContainerClient struct {
	client
	t          *testing.T
	input      mockInput
	response   azcosmos.ItemResponse
	err        error
	funcCalled bool
}

func (m *mockCosmosContainerClient) CreateItem(ctx context.Context, partitionKey azcosmos.PartitionKey, item []byte, o *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
	m.funcCalled = true
	require.Equal(m.t, m.input.partitionKey, partitionKey)
	// require.Equal(m.t, m.input.item, item)
	require.Equal(m.t, m.input.o, o)
	return m.response, m.err
}

func (m *mockCosmosContainerClient) ReadItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, o *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
	return azcosmos.ItemResponse{}, nil
}

func (m *mockCosmosContainerClient) ReplaceItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, item []byte, o *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
	return azcosmos.ItemResponse{}, nil
}

func (m *mockCosmosContainerClient) DeleteItem(ctx context.Context, partitionKey azcosmos.PartitionKey, itemId string, o *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
	return azcosmos.ItemResponse{}, nil
}

func Test_CreateNote_Success(t *testing.T) {
	// Arrange
	mockClient := mockCosmosContainerClient{
		t: t,
		input: mockInput{
			partitionKey: azcosmos.NewPartitionKeyString("category"),
			item:         []byte("{\"id\":\"id\",\"category\":\"category\",\"note\":\"note\"}"),
			o: &azcosmos.ItemOptions{
				EnableContentResponseOnWrite: true,
			},
		},
		response: azcosmos.ItemResponse{
			Value: []byte(`{"test":"test"}`),
		},
		err: nil,
	}

	cosmosDB, err := NewCosmosDB(
		&mockClient,
		log.New(),
	)
	assert.NoError(t, err)

	note := Note{
		ID:       "id",
		Category: "category",
		Note:     "note",
	}

	// Act
	_, err = cosmosDB.CreateNote(context.Background(), note)

	// Assert
	require.NoError(t, err)
	require.True(t, mockClient.funcCalled)
}

func Test_CreateNote_Err_on_Create(t *testing.T) {
	// Arrange
	mockClient := mockCosmosContainerClient{
		t: t,
		input: mockInput{
			partitionKey: azcosmos.NewPartitionKeyString("category"),
			item:         []byte("{\"id\":\"id\",\"category\":\"category\",\"note\":\"note\"}"),
			o: &azcosmos.ItemOptions{
				EnableContentResponseOnWrite: true,
			},
		},
		response: azcosmos.ItemResponse{},
		err:      assert.AnError,
	}

	cosmosDB, err := NewCosmosDB(
		&mockClient,
		log.New(),
	)
	assert.NoError(t, err)

	note := Note{
		ID:       "id",
		Category: "category",
		Note:     "note",
	}

	// Act
	_, err = cosmosDB.CreateNote(context.Background(), note)

	// Assert
	require.Error(t, err)
	require.True(t, mockClient.funcCalled)
}

func Test_CreateNote_Err_on_Unmarshal(t *testing.T) {
	// Arrange
	mockClient := mockCosmosContainerClient{
		t: t,
		input: mockInput{
			partitionKey: azcosmos.NewPartitionKeyString("category"),
			item:         []byte("{\"id\":\"id\",\"category\":\"category\",\"note\":\"note\"}"),
			o: &azcosmos.ItemOptions{
				EnableContentResponseOnWrite: true,
			},
		},
		response: azcosmos.ItemResponse{
			Value: []byte(`notajson`),
		},
		err: nil,
	}

	cosmosDB, err := NewCosmosDB(
		&mockClient,
		log.New(),
	)
	assert.NoError(t, err)

	note := Note{
		ID:       "id",
		Category: "category",
		Note:     "note",
	}

	// Act
	_, err = cosmosDB.CreateNote(context.Background(), note)

	// Assert
	require.Error(t, err)
	require.True(t, mockClient.funcCalled)
}
