package repository

import (
	"randomInsultService/datasource"
	"randomInsultService/mocks"
	"randomInsultService/model"
	"testing"
)

func TestFireStore_InsertEntry(t *testing.T) {
	tests := []struct {
		name string
		fs   fireStore
		ds   datasource.DataSource
		in   struct {
			insultContent model.InsultContent
		}
		expected struct {
			id  *string
			err error
		}
	}{
		{
			name: "error initializing data source",
			fs:   fireStore{},
			ds:   &mocks.DataSource{},
			in: struct {
				insultContent model.InsultContent
			}{
				insultContent: model.InsultContent{
					Verb:      "test verb",
					Adjective: "test adjective",
					Noun:      "test noun ",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
