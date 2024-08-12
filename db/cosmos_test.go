package db

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// What we want to test is that your CreateNote method works as expected.
// By creating an ID for the Note, and that it handles responses from
// the CosmosDB client correctly.
func TestCosmosDB_CreateNote(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			client client
			note   Note
		}
		want    Note
		wantErr error
	}{
		{
			name: "Create note",
			input: struct {
				client client
				note   Note
			}{
				client: &mockCosmosClient{},
				note: Note{
					Category: "category",
					Note:     "note",
				},
			},
			want: Note{
				ID:       "id",
				Category: "category",
				Note:     "note",
			},
		},
	}

	newUUID = func() string {
		return "id"
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := &CosmosDB{
				cl: test.input.client,
			}

			got, gotErr := c.CreateNote(context.Background(), test.input.note)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateNote() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr); diff != "" {
				t.Errorf("CreateNote() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

type mockCosmosClient struct {
	data map[string][]byte
	err  error
}

// The mock client implementation of CreateItem stores the input item in a map (after unmarshalling it to retreive the ID for
// the key) and returns the input item as the response just as the CosmosDB client would.
func (c *mockCosmosClient) CreateItem(ctx context.Context, partitionKey string, item []byte) ([]byte, error) {
	if c.err != nil {
		return []byte{}, c.err
	}
	if c.data == nil {
		c.data = make(map[string][]byte)
	}

	var note Note
	if err := json.Unmarshal(item, &note); err != nil {
		return []byte{}, err
	}

	c.data[note.ID] = item
	return item, nil
}

func (c mockCosmosClient) ReplaceItem(ctx context.Context, partitionKey, id string, item []byte) ([]byte, error) {
	if c.err != nil {
		return []byte{}, c.err
	}
	return []byte{}, nil
}

func (c mockCosmosClient) DeleteItem(ctx context.Context, partitionKey, id string) error {
	if c.err != nil {
		return c.err
	}
	return nil
}

func (c mockCosmosClient) ReadItem(ctx context.Context, partitionKey, id string) ([]byte, error) {
	if c.err != nil {
		return []byte{}, c.err
	}
	return []byte{}, nil
}

func (c mockCosmosClient) ListItems(ctx context.Context, partitionKey string) ([][]byte, error) {
	if c.err != nil {
		return [][]byte{}, c.err
	}
	return [][]byte{}, nil
}
