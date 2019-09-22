package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/juju/loggo"
	"randomInsultService/datasource"
	"randomInsultService/model"
)

// Firebase is an interface that contains methods pertaining to CRUD operations for FireStore db
type FireStore interface {
	InsertEntry(insultContent model.InsultContent) (*string, error)
}

type fireStore struct {
	ctx context.Context
	log loggo.Logger
	ds  datasource.DataSource
}

// NewFireBase creates a new firestore repository
func NewFireBase(ctx context.Context, ds datasource.DataSource, log loggo.Logger) FireStore {
	return &fireStore{
		ctx: ctx,
		ds:  ds,
		log: log,
	}
}

// InsertEntry will insert a generated insult into the firestore DB and return a value corresponding to its ID
//  If unsuccessful, it will bubble up an error
func (f *fireStore) InsertEntry(insultContent model.InsultContent) (*string, error) {
	client, err := f.ds.InitializeDataSource()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error Intializing data source: %v", err))
	}
	defer client.Close()
	insults := client.Collection("insults")
	//Should insert generated insult into firestore
	d, wr, err := insults.Add(f.ctx, insultContent)
	if err != nil {
		//Should allow error to bubble up upon failure
		f.log.Errorf("failed to create an insult doc for write: %v", err)
		return nil, err
	}
	fmt.Printf("documentRef: %v\n", d)
	fmt.Printf("writeResult: %v\n", wr)

	return &d.ID, nil
}
