package datastore

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

// Client wraps the Google Cloud Datastore client
type Client struct {
	client *datastore.Client
}

// NewClient creates a new Datastore client
func NewClient(ctx context.Context, projectID string) (*Client, error) {
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create datastore client: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

// Close closes the Datastore client
func (c *Client) Close() error {
	return c.client.Close()
}

// GetClient returns the underlying Datastore client
func (c *Client) GetClient() *datastore.Client {
	return c.client
}
