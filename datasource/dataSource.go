package datasource

import (
	. "cloud.google.com/go/firestore"
	"context"
	"github.com/juju/loggo"
	"google.golang.org/api/option"
	"randomInsultService/model"
)

// DataSource defines what functions we need from the database
type DataSource interface {
	OpenConnection() error
	CloseConnection() error
	ReadEntries() (*DocumentSnapshot, error)
	InsertEntry(content model.InsultContent) (*DocumentRef, *WriteResult, error)
}

type fireStoreDB struct {
	log       loggo.Logger
	ctx       context.Context
	projectId string
	client    *Client
}

// NewDataSource returns a Datasource interface
func NewDataSource(l loggo.Logger, ctx context.Context, projectId string) DataSource {
	return &fireStoreDB{
		log:       l,
		ctx:       ctx,
		projectId: projectId,
		client:    nil,
	}
}

// OpenConnection returns an error if one is present.  It also sets a client within the fireStoreDB object.
// The error is returned when firestore.NewClient fails to create a client.
func (f *fireStoreDB) OpenConnection() error {
	const jsonPath = "../insult-41aadfdfb47b.json"
	client, err := NewClient(f.ctx, f.projectId, option.WithCredentialsFile(jsonPath))
	if err != nil {
		f.log.Errorf("error initializing Fire Store client with projectId: %s. Received error: %v", f.projectId, err)
		return err
	}
	f.client = client
	return nil
}

func (f *fireStoreDB) InsertEntry(content model.InsultContent) (*DocumentRef, *WriteResult, error) {
	insults := f.client.Collection("insults")
	//Should insert generated insult into firestore
	d, wr, err := insults.Add(f.ctx, content)
	if err != nil {
		//Should allow error to bubble up upon failure
		f.log.Errorf("failed to create an insult doc with content: %v for write: %v", content, err)
		return nil, nil, err
	}
	return d, wr, nil
}
func (f *fireStoreDB) ReadEntries() (*DocumentSnapshot, error) {
	words, err := f.client.Collection("insults").Doc("words").Get(f.ctx)
	if err != nil {
		f.log.Errorf("failed to read possible words with error: %v", err)
		return nil, err
	}
	return words, nil
}
func (f *fireStoreDB) CloseConnection() error {
	err := f.client.Close()
	if err != nil {
		f.log.Errorf("error when closing the connection: %v", err)
		return err
	}
	f.client = nil
	return nil
}
