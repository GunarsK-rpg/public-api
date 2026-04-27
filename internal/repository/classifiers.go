package repository

import (
	"context"
	"encoding/json"
)

// ClassifierRepository defines methods for classifier data access.
// Scope-level access checks also live inside classifiers.get_all_classifiers
// so direct DB callers are protected; the Go wrappers here exist so handlers
// can gate the Redis cache before hitting it.
type ClassifierRepository interface {
	GetClassifiersFiltered(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error)
	GetSourceBooks(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetSourceBookDependencyIDs(ctx context.Context, auth AuthContext, sourceBookID int64) ([]int64, error)
	RequireSourceBookAccessible(ctx context.Context, auth AuthContext, sourceBookID int64) error
	ValidateHeroAccess(ctx context.Context, auth AuthContext, heroID int64) error
}

func (r *repository) GetClassifiersFiltered(ctx context.Context, auth AuthContext, filter json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_all_classifiers($1::jsonb)", filter)
}

func (r *repository) GetSourceBooks(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, `SELECT classifiers.get_source_books('{"includeGlobal": true}'::jsonb)`)
}

func (r *repository) GetSourceBookDependencyIDs(ctx context.Context, auth AuthContext, sourceBookID int64) ([]int64, error) {
	raw, err := r.callFunc(ctx, auth, "SELECT to_jsonb(classifiers.get_source_book_dependency_ids($1))", sourceBookID)
	if err != nil {
		return nil, err
	}
	var ids []int64
	if err := json.Unmarshal(raw, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *repository) RequireSourceBookAccessible(ctx context.Context, auth AuthContext, sourceBookID int64) error {
	_, err := r.callFunc(ctx, auth, "SELECT public.require_source_book_accessible($1)", sourceBookID)
	return err
}

func (r *repository) ValidateHeroAccess(ctx context.Context, auth AuthContext, heroID int64) error {
	_, err := r.callFunc(ctx, auth, "SELECT public.require_hero_access($1)", heroID)
	return err
}
