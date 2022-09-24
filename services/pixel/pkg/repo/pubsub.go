package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/mikestefanello/gcpworkflow/config"
	"github.com/mikestefanello/gcpworkflow/pkg/pixel"
)

type PubSubRepo struct {
	client *pubsub.Client
	topic  string
}

func NewPubSubRepo(cfg config.Config) (pixel.Repository, error) {
	client, err := pubsub.NewClient(context.Background(), cfg.Cloud.Project)
	if err != nil {
		return nil, fmt.Errorf("pubsub: NewClient: %v", err)
	}

	return &PubSubRepo{
		client: client,
		topic:  cfg.Cloud.PubSubTopic,
	}, nil
}

func (r *PubSubRepo) Store(ctx context.Context, p pixel.Pixel) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	t := r.client.Topic(r.topic)
	result := t.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("pubsub: result.Get: %v", err)
	}

	return nil
}

func (r *PubSubRepo) Close() error {
	return r.client.Close()
}
