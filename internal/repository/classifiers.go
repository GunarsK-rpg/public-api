package repository

import (
	"context"
	"encoding/json"
)

// ClassifierRepository defines methods for classifier data access.
type ClassifierRepository interface {
	GetAllClassifiers(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetSourceBooks(ctx context.Context, auth AuthContext) (json.RawMessage, error)
}

func (r *repository) GetAllClassifiers(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_all_classifiers()")
}

func (r *repository) GetSourceBooks(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_source_books()")
}
