package firestore

import (
	"fmt"

	gcfs "cloud.google.com/go/firestore"
	"github.com/parish/internal/domain"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func isNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}

func scanDocuments[T domain.Entity[T]](iter *gcfs.DocumentIterator) ([]*T, error) {
	var out []*T

	name := (*new(T)).EntityKind()

	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to scan %s documents: %w", name, err)
		}
		var t T
		if err := doc.DataTo(&t); err != nil {
			return nil, fmt.Errorf("failed to decode % document: %w", name, err)
		}

		t = t.SetID(doc.Ref.ID)
		out = append(out, &t)
	}

	return out, nil
}
