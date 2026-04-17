package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/GunarsK-rpg/public-api/internal/constants"
)

// HomebrewRepository defines methods for homebrew source-book and classifier writes.
// All methods are thin wrappers around classifiers.* SQL functions. Ownership and
// validation are enforced at the DB layer.
type HomebrewRepository interface {
	// Source book CRUD (keyed by UUID code, matching the existing DB function shapes)
	UpsertSourceBook(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	GetSourceBookByCode(ctx context.Context, auth AuthContext, code string) (json.RawMessage, error)
	DeleteSourceBookByCode(ctx context.Context, auth AuthContext, code string) (bool, error)
	// RestoreSourceBookByCode returns the restored book as JSONB (matching the
	// upsert_* return-the-row convention). Returns NULL JSONB if already active,
	// not found, or not owned; handler maps that to 404.
	RestoreSourceBookByCode(ctx context.Context, auth AuthContext, code string) (json.RawMessage, error)

	// ListMyHomebrewSourceBooks returns the session user's own homebrew books,
	// including inactive and soft-deleted rows. Used by the Library page.
	ListMyHomebrewSourceBooks(ctx context.Context, auth AuthContext) (json.RawMessage, error)

	// Generic classifier CRUD. classifierType is the URL plural form;
	// the repository validates it against the allow-list before any string
	// interpolation into SQL.
	UpsertClassifier(ctx context.Context, auth AuthContext, classifierType string, data json.RawMessage) (json.RawMessage, error)
	DeleteClassifier(ctx context.Context, auth AuthContext, classifierType string, id int64) (bool, error)
	// RestoreClassifier returns the restored row as JSONB via its detail view.
	// NULL JSONB means already active or not found; handler maps to 404.
	RestoreClassifier(ctx context.Context, auth AuthContext, classifierType string, id int64) (json.RawMessage, error)

	// IsClassifierInScope reports whether a classifier row's scope matches
	// (sourceBookID, heroID). Either bound may be nil; nil matches nil.
	// Propagates no_data_found when the row does not exist.
	IsClassifierInScope(ctx context.Context, auth AuthContext, classifierType string, id int64, sourceBookID, heroID *int64) (bool, error)
}

// ErrUnknownClassifierType is returned when a caller passes a classifier type
// that is not in the allow-list. Belt-and-suspenders alongside handler-level
// validation: the repository is the last guard before SQL string interpolation.
var ErrUnknownClassifierType = fmt.Errorf("unknown classifier type")

func (r *repository) UpsertSourceBook(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.upsert_source_book($1::jsonb)", data)
}

func (r *repository) GetSourceBookByCode(ctx context.Context, auth AuthContext, code string) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.get_source_book_by_code($1::uuid)", code)
}

func (r *repository) DeleteSourceBookByCode(ctx context.Context, auth AuthContext, code string) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT classifiers.delete_source_book($1::uuid)", code)
}

func (r *repository) RestoreSourceBookByCode(ctx context.Context, auth AuthContext, code string) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT classifiers.restore_source_book($1::uuid)", code)
}

func (r *repository) ListMyHomebrewSourceBooks(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth,
		`SELECT classifiers.get_source_books('{"includeInactive": true, "includeDeleted": true}'::jsonb)`)
}

func (r *repository) UpsertClassifier(ctx context.Context, auth AuthContext, classifierType string, data json.RawMessage) (json.RawMessage, error) {
	suffix, ok := constants.ClassifierTypeSuffix(classifierType)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownClassifierType, classifierType)
	}
	query := fmt.Sprintf("SELECT classifiers.upsert_%s($1::jsonb)", suffix)
	return r.callFunc(ctx, auth, query, data)
}

func (r *repository) DeleteClassifier(ctx context.Context, auth AuthContext, classifierType string, id int64) (bool, error) {
	suffix, ok := constants.ClassifierTypeSuffix(classifierType)
	if !ok {
		return false, fmt.Errorf("%w: %q", ErrUnknownClassifierType, classifierType)
	}
	query := fmt.Sprintf("SELECT classifiers.delete_%s($1)", suffix)
	return r.execFunc(ctx, auth, query, id)
}

func (r *repository) RestoreClassifier(ctx context.Context, auth AuthContext, classifierType string, id int64) (json.RawMessage, error) {
	table, ok := constants.ClassifierTableName(classifierType)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownClassifierType, classifierType)
	}
	return r.callFunc(ctx, auth, "SELECT classifiers.restore_classifier($1, $2)", table, id)
}

func (r *repository) IsClassifierInScope(ctx context.Context, auth AuthContext, classifierType string, id int64, sourceBookID, heroID *int64) (bool, error) {
	table, ok := constants.ClassifierTableName(classifierType)
	if !ok {
		return false, fmt.Errorf("%w: %q", ErrUnknownClassifierType, classifierType)
	}
	return r.execFunc(ctx, auth,
		"SELECT public.is_classifier_in_scope($1, $2, $3, $4)",
		table, id, sourceBookID, heroID,
	)
}
