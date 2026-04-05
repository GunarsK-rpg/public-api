package repository

import (
	"context"
	"encoding/json"
)

// ClassifierRepository defines methods for classifier data access.
type ClassifierRepository interface {
	GetClassifiersFiltered(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error)
	GetSourceBooks(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	ValidateHeroAccess(ctx context.Context, auth AuthContext, heroID int64) error
}

func (r *repository) GetClassifiersFiltered(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_all_classifiers($1::jsonb)", filter)
}

func (r *repository) GetSourceBooks(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_source_books()")
}

func (r *repository) ValidateHeroAccess(ctx context.Context, auth AuthContext, heroID int64) error {
	_, err := r.callFunc(ctx, auth, "SELECT public.require_hero_access($1)", heroID)
	return err
}
