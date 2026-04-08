package firestore

import (
	"context"
	"fmt"

	gcfs "cloud.google.com/go/firestore"
)

// Store wraps the Google Cloud Firestore client (native API).
type Store struct {
	*gcfs.Client
}

// NewStore creates a Firestore client. Uses FIRESTORE_EMULATOR_HOST when set (local dev).
func NewStore(ctx context.Context, projectID string) (*Store, error) {
	client, err := gcfs.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}
	return &Store{Client: client}, nil
}

// Close closes the Firestore client.
func (s *Store) Close() error {
	return s.Client.Close()
}
