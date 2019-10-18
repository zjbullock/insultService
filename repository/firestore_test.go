package repository

import (
	"cloud.google.com/go/firestore"
	"errors"
	"fmt"
	"github.com/go-test/deep"
	"github.com/juju/loggo"
	"github.com/stretchr/testify/assert"
	"insultService/datasource"
	"insultService/mocks"
	"insultService/model"
	"testing"
)

func TestNewFireStore(t *testing.T) {
	tests := []struct {
		name string
		in   struct {
			ds  datasource.DataSource
			log loggo.Logger
		}
		expected FireStore
	}{
		{
			in: struct {
				ds  datasource.DataSource
				log loggo.Logger
			}{
				ds:  nil,
				log: loggo.Logger{},
			},
			expected: &fireStore{
				ds:  nil,
				log: loggo.Logger{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testCase := assert.New(t)
			res := NewFireStore(test.in.ds, test.in.log)
			testCase.Equal(test.expected, res)
		})
	}
}

func TestFireStore_InsertEntry(t *testing.T) {
	tests := []struct {
		name string
		fs   fireStore
		ds   struct {
			openConnection  error
			closeConnection error
			insertEntry     struct {
				doc *firestore.DocumentRef
				res *firestore.WriteResult
				err error
			}
		}
		in struct {
			insultContent model.InsultContent
		}
		expected struct {
			id  *string
			err error
		}
	}{
		{
			name: "error opening datasource connection",
			fs: fireStore{
				log: loggo.Logger{},
			},
			ds: struct {
				openConnection  error
				closeConnection error
				insertEntry     struct {
					doc *firestore.DocumentRef
					res *firestore.WriteResult
					err error
				}
			}{
				openConnection:  errors.New("error opening the connection"),
				closeConnection: nil,
				insertEntry: struct {
					doc *firestore.DocumentRef
					res *firestore.WriteResult
					err error
				}{
					doc: nil,
					res: nil,
					err: nil,
				},
			},
			in: struct {
				insultContent model.InsultContent
			}{
				insultContent: model.InsultContent{
					Verb:      "test verb",
					Adjective: "test adjective",
					Noun:      "test noun ",
				},
			},
			expected: struct {
				id  *string
				err error
			}{
				id:  nil,
				err: errors.New(fmt.Sprintf("error Intializing data source: %v", errors.New("error opening the connection"))),
			},
		},
		{
			name: "error Inserting an entry into the DB",
			fs: fireStore{
				log: loggo.Logger{},
			},
			ds: struct {
				openConnection  error
				closeConnection error
				insertEntry     struct {
					doc *firestore.DocumentRef
					res *firestore.WriteResult
					err error
				}
			}{
				openConnection:  nil,
				closeConnection: nil,
				insertEntry: struct {
					doc *firestore.DocumentRef
					res *firestore.WriteResult
					err error
				}{
					doc: nil,
					res: nil,
					err: errors.New("error inserting doc"),
				},
			},
			in: struct {
				insultContent model.InsultContent
			}{
				insultContent: model.InsultContent{
					Verb:      "test verb",
					Adjective: "test adjective",
					Noun:      "test noun ",
				},
			},
			expected: struct {
				id  *string
				err error
			}{
				id: nil,
				err: errors.New(fmt.Sprintf("error inserting insult doc: %v.  received error: %v", model.InsultContent{
					Verb:      "test verb",
					Adjective: "test adjective",
					Noun:      "test noun ",
				}, errors.New("error inserting doc"))),
			},
		},
		{
			name: "successfully inserted document into DB and retrieved its ID",
			fs: fireStore{
				log: loggo.Logger{},
			},
			ds: struct {
				openConnection  error
				closeConnection error
				insertEntry     struct {
					doc *firestore.DocumentRef
					res *firestore.WriteResult
					err error
				}
			}{
				openConnection:  nil,
				closeConnection: nil,
				insertEntry: struct {
					doc *firestore.DocumentRef
					res *firestore.WriteResult
					err error
				}{
					doc: &firestore.DocumentRef{
						ID: "test id",
					},
					res: &firestore.WriteResult{},
					err: nil,
				},
			},
			in: struct {
				insultContent model.InsultContent
			}{
				insultContent: model.InsultContent{
					Verb:      "test verb",
					Adjective: "test adjective",
					Noun:      "test noun ",
				},
			},
			expected: struct {
				id  *string
				err error
			}{
				id: func() *string {
					s := "test id"
					return &s
				}(),
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testCase := assert.New(t)
			ds := &mocks.DataSource{}
			ds.On("OpenConnection").Return(test.ds.openConnection)
			ds.On("CloseConnection").Return(test.ds.closeConnection)
			ds.On("InsertEntry", test.in.insultContent).Return(test.ds.insertEntry.doc, test.ds.insertEntry.res, test.ds.insertEntry.err)
			test.fs.ds = ds
			id, err := test.fs.InsertEntry(test.in.insultContent)
			testCase.Equal(test.expected.id, id)
			testCase.Len(deep.Equal(test.expected.err, err), 0)

		})
	}
}
