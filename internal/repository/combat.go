package repository

import (
	"context"
	"encoding/json"
	"fmt"
)

// CombatRepository defines methods for combat data access.
type CombatRepository interface {
	// NPCs
	GetNpcOptions(ctx context.Context, auth AuthContext, campaignID int64) (json.RawMessage, error)
	GetNpc(ctx context.Context, auth AuthContext, id int64, campaignID int64) (json.RawMessage, error)
	UpsertNpc(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	DeleteNpc(ctx context.Context, auth AuthContext, id int64, campaignID int64) (bool, error)

	// Combats
	GetCombats(ctx context.Context, auth AuthContext, campaignID int64) (json.RawMessage, error)
	GetCombat(ctx context.Context, auth AuthContext, id int64, campaignID int64) (json.RawMessage, error)
	UpsertCombat(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	DeleteCombat(ctx context.Context, auth AuthContext, id int64, campaignID int64) (bool, error)

	// Combat NPC instances
	UpsertCombatNpc(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	DeleteCombatNpc(ctx context.Context, auth AuthContext, id int64, combatID int64, campaignID int64) (bool, error)

	// Combat NPC resource patches
	PatchCombatNpcHp(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	PatchCombatNpcFocus(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	PatchCombatNpcInvestiture(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
}

// NPCs

func (r *repository) GetNpcOptions(ctx context.Context, auth AuthContext, campaignID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.get_npc_options($1)", campaignID)
}

func (r *repository) GetNpc(ctx context.Context, auth AuthContext, id int64, campaignID int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.get_npc($1, $2)", id, campaignID)
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

// Combat NPC instances

func (r *repository) UpsertCombatNpc(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	campaignID, err := extractInt64(data, "campaignId")
	if err != nil {
		return nil, err
	}
	combatID, err := extractInt64(data, "combatId")
	if err != nil {
		return nil, err
	}
	return r.callFunc(ctx, auth, "SELECT combat.upsert_combat_npc($1, $2, $3::jsonb)", campaignID, combatID, data)
}

func (r *repository) DeleteCombatNpc(ctx context.Context, auth AuthContext, id int64, combatID int64, campaignID int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT combat.delete_combat_npc($1, $2, $3)", id, combatID, campaignID)
}

// Combat NPC resource patches

func (r *repository) PatchCombatNpcHp(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.patch_combat_npc_hp($1::jsonb)", data)
}

func (r *repository) PatchCombatNpcFocus(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.patch_combat_npc_focus($1::jsonb)", data)
}

func (r *repository) PatchCombatNpcInvestiture(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT combat.patch_combat_npc_investiture($1::jsonb)", data)
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
