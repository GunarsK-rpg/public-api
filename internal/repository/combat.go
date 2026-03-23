package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// CombatRepository defines methods for combat data access.
type CombatRepository interface {
	// NPCs (templates)
	GetNpcOptions(ctx context.Context, auth AuthContext, campaignID int64) (json.RawMessage, error)
	GetNpcLibrary(ctx context.Context, auth AuthContext, campaignID int64) (json.RawMessage, error)
	GetNpc(ctx context.Context, auth AuthContext, id int64, campaignID int64) (json.RawMessage, error)
	GetNpcById(ctx context.Context, auth AuthContext, id int64) (json.RawMessage, error)
	UpsertNpc(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	DeleteNpc(ctx context.Context, auth AuthContext, id int64, campaignID int64) (bool, error)
	UpsertNpcAvatar(ctx context.Context, auth AuthContext, npcID int64, campaignID int64, avatarKey string) error
	DeleteNpcAvatar(ctx context.Context, auth AuthContext, npcID int64, campaignID int64) error

	// Combats
	GetCombats(ctx context.Context, auth AuthContext, campaignID int64) (json.RawMessage, error)
	GetCombat(ctx context.Context, auth AuthContext, id int64, campaignID int64) (json.RawMessage, error)
	UpsertCombat(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	DeleteCombat(ctx context.Context, auth AuthContext, id int64, campaignID int64) (bool, error)

	// Combat round management
	EndCombatRound(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)

	// NPC instances (combat + companion)
	GetNpcInstance(ctx context.Context, auth AuthContext, id int64) (json.RawMessage, error)
	CreateNpcInstance(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	PatchNpcInstance(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	DeleteNpcInstance(ctx context.Context, auth AuthContext, id int64) (bool, error)

	// Companion-specific queries
	GetHeroNpcInstances(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)
	GetCompanionNpcOptions(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error)
}

// NPCs (templates)

func (r *repository) GetNpcOptions(ctx context.Context, auth AuthContext, campaignID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.get_npc_options($1)", campaignID)
}

func (r *repository) GetNpcLibrary(ctx context.Context, auth AuthContext, campaignID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.get_npc_library($1)", campaignID)
}

func (r *repository) GetNpc(ctx context.Context, auth AuthContext, id int64, campaignID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.get_npc($1, $2)", id, campaignID)
}

func (r *repository) GetNpcById(ctx context.Context, auth AuthContext, id int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.get_npc($1)", id)
}

func (r *repository) UpsertNpc(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	campaignID, err := extractInt64(data, "campaignId")
	if err != nil {
		return nil, err
	}
	return r.callFunc(ctx, auth, "SELECT combat.upsert_npc($1, $2::jsonb)", campaignID, data)
}

func (r *repository) DeleteNpc(ctx context.Context, auth AuthContext, id int64, campaignID int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT combat.delete_npc($1, $2)", id, campaignID)
}

func (r *repository) UpsertNpcAvatar(ctx context.Context, auth AuthContext, npcID int64, campaignID int64, avatarKey string) error {
	return r.withAuditTx(ctx, auth, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "SELECT combat.upsert_npc_avatar($1, $2, $3)", npcID, campaignID, avatarKey)
		return err
	})
}

func (r *repository) DeleteNpcAvatar(ctx context.Context, auth AuthContext, npcID int64, campaignID int64) error {
	return r.withAuditTx(ctx, auth, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "SELECT combat.delete_npc_avatar($1, $2)", npcID, campaignID)
		return err
	})
}

// Combats

func (r *repository) GetCombats(ctx context.Context, auth AuthContext, campaignID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.get_combats($1)", campaignID)
}

func (r *repository) GetCombat(ctx context.Context, auth AuthContext, id int64, campaignID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.get_combat($1, $2)", id, campaignID)
}

func (r *repository) UpsertCombat(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	campaignID, err := extractInt64(data, "campaignId")
	if err != nil {
		return nil, err
	}
	return r.callFunc(ctx, auth, "SELECT combat.upsert_combat($1, $2::jsonb)", campaignID, data)
}

func (r *repository) DeleteCombat(ctx context.Context, auth AuthContext, id int64, campaignID int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT combat.delete_combat($1, $2)", id, campaignID)
}

// Combat round management

func (r *repository) EndCombatRound(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.end_combat_round($1::jsonb)", data)
}

// NPC instances (combat + companion)

func (r *repository) GetNpcInstance(ctx context.Context, auth AuthContext, id int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.get_npc_instance($1)", id)
}

func (r *repository) CreateNpcInstance(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.upsert_npc_instance($1::jsonb)", data)
}

func (r *repository) PatchNpcInstance(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	// Route to resource patch if "field" key is present, otherwise metadata update
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	if _, hasField := payload["field"]; hasField {
		return r.callFunc(ctx, auth, "SELECT combat.patch_npc_instance_resource($1::jsonb)", data)
	}
	return r.callFunc(ctx, auth, "SELECT combat.upsert_npc_instance($1::jsonb)", data)
}

func (r *repository) DeleteNpcInstance(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT combat.delete_npc_instance($1)", id)
}

// Companion-specific queries

func (r *repository) GetHeroNpcInstances(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.get_hero_npc_instances($1)", heroID)
}

func (r *repository) GetCompanionNpcOptions(ctx context.Context, auth AuthContext, heroID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.get_companion_npc_options($1)", heroID)
}

// extractInt64 pulls a required integer field from JSON data.
func extractInt64(data json.RawMessage, field string) (int64, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return 0, fmt.Errorf("invalid JSON: %w", err)
	}
	raw, ok := m[field]
	if !ok {
		return 0, fmt.Errorf("%s is required", field)
	}
	var val int64
	if err := json.Unmarshal(raw, &val); err != nil {
		return 0, fmt.Errorf("invalid %s: must be an integer", field)
	}
	return val, nil
}
