package repository

import (
	"errors"
	"fmt"
	"github.com/juju/loggo"
	"google.golang.org/api/iterator"
	"insultService/datasource"
	"insultService/model"
)

// Firebase is an interface that contains methods pertaining to CRUD operations for FireStore db
type FireStore interface {
	ReadAllWords() (*model.Words, error)
	InsertInsultEntry(insultContent model.InsultContent) (*string, error)
	ReadInsults() ([]model.InsultContent, error)
}

type fireStore struct {
	log loggo.Logger
	ds  datasource.DataSource
}

// NewFireStore creates a new firestore repository
func NewFireStore(ds datasource.DataSource, log loggo.Logger) FireStore {
	return &fireStore{
		ds:  ds,
		log: log,
	}
}

func (f *fireStore) ReadAllWords() (*model.Words, error) {
	err := f.ds.OpenConnection()
	if err != nil {
		f.log.Errorf("error initializing datasource: %v", err)
		return nil, errors.New(fmt.Sprintf("error initializing data source: %v", err))
	}
	defer f.ds.CloseConnection()
	doc, err := f.ds.ReadEntries()
	if err != nil {
		f.log.Errorf("error reading from datasource: %v", err)
		return nil, errors.New(fmt.Sprintf("error reading from datasource: %v", err))
	}
	var words model.Words
	if err := doc.DataTo(&words); err != nil {
		f.log.Errorf("error converting document snap to a Words model")
		return nil, errors.New(fmt.Sprintf("error converting document snap to a model.Words, :%v", err))
	}
	return &words, nil
}

// InsertInsultEntry will insert a generated insult into the firestore DB and return a value corresponding to its ID
//  If unsuccessful, it will bubble up an error
func (f *fireStore) InsertInsultEntry(insultContent model.InsultContent) (*string, error) {
	err := f.ds.OpenConnection()
	if err != nil {
		f.log.Errorf("error initializing datasource: %v", err)
		return nil, errors.New(fmt.Sprintf("error Intializing data source: %v", err))
	}
	defer f.ds.CloseConnection()
	d, _, err := f.ds.InsertEntry(insultContent, "insult")
	if err != nil {
		//Should allow error to bubble up upon failure
		f.log.Errorf("failed to create an insult doc for write: %v", err)
		return nil, errors.New(fmt.Sprintf("error inserting insult doc: %v.  received error: %v", insultContent, err))
	}
	return &d.ID, nil
}

func (f *fireStore) ReadInsults() (insultContents []model.InsultContent, err error) {
	err = f.ds.OpenConnection()
	if err != nil {
		f.log.Errorf("error initializing datasource: %v", err)
		return nil, errors.New(fmt.Sprintf("error Intializing data source: %v", err))
	}
	defer f.ds.CloseConnection()
	d := f.ds.ReadCollection("insults")
	defer d.Stop()
	for {
		doc, err := d.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			f.log.Errorf("failed to read document insult contents: %v", err)
		}
		var insultContent model.InsultContent
		if err := doc.DataTo(&insultContent); err != nil {
			f.log.Errorf("error converting document snap to an insult content model")
			return nil, errors.New(fmt.Sprintf("error converting document snap to a model.InsultContent, :%v", err))
		}
		insultContents = append(insultContents, insultContent)
	}
	return insultContents, nil
}
