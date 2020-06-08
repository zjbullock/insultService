package datasource

import (
	. "cloud.google.com/go/firestore"
	"context"
	"github.com/juju/loggo"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"insultService/model"
)

// DataSource defines what functions we need from the database
type DataSource interface {
	OpenConnection() error
	CloseConnection() error
	ReadEntries(collection, doc string) (*DocumentSnapshot, error)
	InsertEntry(content interface{}, collection string) (*DocumentRef, *WriteResult, error)
	ReadCollection(collection string) *DocumentIterator
	QueryCollection(collection string, args []model.QueryArg) ([]*DocumentSnapshot, error)
	UpdateDocument(collection, doc string, content interface{}) error
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
	const jsonPath = "./datasource/insultDB.json"
	client, err := NewClient(f.ctx, f.projectId, option.WithCredentialsFile(jsonPath))
	if err != nil {
		f.log.Errorf("error initializing Fire Store client with projectId: %s. Received error: %v", f.projectId, err)
		return err
	}
	f.client = client
	return nil
}

func (f *fireStoreDB) InsertEntry(content interface{}, collection string) (*DocumentRef, *WriteResult, error) {
	insults := f.client.Collection(collection)
	//Should insert generated insult into firestore
	d, wr, err := insults.Add(f.ctx, content)
	if err != nil {
		//Should allow error to bubble up upon failure
		f.log.Errorf("failed to create an insult doc with content: %v for write: %v", content, err)
		return nil, nil, err
	}
	return d, wr, nil
}

func (f *fireStoreDB) ReadEntries(collection, doc string) (*DocumentSnapshot, error) {
	words, err := f.client.Collection(collection).Doc(doc).Get(f.ctx)
	if err != nil {
		f.log.Errorf("failed to read possible words with error: %v", err)
		return nil, err
	}
	return words, nil
}

func (f *fireStoreDB) ReadCollection(collection string) *DocumentIterator {
	insultContent := f.client.Collection(collection).Documents(f.ctx)
	return insultContent
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

func (f *fireStoreDB) QueryCollection(collection string, args []model.QueryArg) ([]*DocumentSnapshot, error) {
	q := f.client.Collection(collection).Query
	for _, arg := range args {
		q = q.Where(arg.Path, arg.Op, arg.Value)
	}
	iter := q.Documents(f.ctx)
	defer iter.Stop()
	var docs []*DocumentSnapshot
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			f.log.Errorf("error iterating through queried documents with error: %v", err)
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func (f *fireStoreDB) UpdateDocument(collection, doc string, content interface{}) error {
	f.log.Infof("collection name: %s", collection)
	col := f.client.Collection(collection)
	f.log.Infof("collection: %v", col)
	//Should insert generated insult into firestore
	_, err := col.Doc(doc).Set(f.ctx, content)
	if err != nil {
		//Should allow error to bubble up upon failure
		f.log.Errorf("failed to update doc %v with content: %v for write: %v", doc, content, err)
		return err
	}
	return nil
}
