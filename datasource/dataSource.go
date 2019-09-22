package datasource

import (
	. "cloud.google.com/go/firestore"
	"context"
	"github.com/juju/loggo"
	"google.golang.org/api/option"
)

// DataSource
type DataSource interface {
	InitializeDataSource() (*Client, error)
}

type fireBase struct {
	log       loggo.Logger
	ctx       context.Context
	projectId string
}

// NewDataSource returns a Datasource interface
func NewDataSource(l loggo.Logger, ctx context.Context, projectId string) DataSource {
	return &fireBase{
		log:       l,
		ctx:       ctx,
		projectId: projectId,
	}
}

// InitializeDataSource returns a pointer to a firestore client an error if one is present.
// The error is returned when firestore.NewClient fails to create a client.
func (f *fireBase) InitializeDataSource() (*Client, error) {
	const jsonPath = "../insult-41aadfdfb47b.json"
	client, err := NewClient(f.ctx, f.projectId, option.WithCredentialsFile(jsonPath))
	if err != nil {
		f.log.Errorf("error initializing Fire Store client with projectId: %s. Received error: %v", f.projectId, err)
		return nil, err
	}
	return client, nil
}
