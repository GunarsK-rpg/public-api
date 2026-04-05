package repository

import (
	"context"
	"encoding/json"
)

// CampaignRepository defines methods for campaign data access.
type CampaignRepository interface {
	GetCampaigns(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetCampaign(ctx context.Context, auth AuthContext, id int64) (json.RawMessage, error)
	GetCampaignByCode(ctx context.Context, auth AuthContext, code string) (json.RawMessage, error)
	UpsertCampaign(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	DeleteCampaign(ctx context.Context, auth AuthContext, id int64) (bool, error)
	RemoveHeroFromCampaign(ctx context.Context, auth AuthContext, heroID int64, campaignID int64) (bool, error)
	GetCampaignSourceBookIDs(ctx context.Context, auth AuthContext, campaignID int64) ([]int64, error)
}

func (r *repository) GetCampaigns(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT campaign.get_campaigns()")
}

func (r *repository) GetCampaign(ctx context.Context, auth AuthContext, id int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT campaign.get_campaign($1)", id)
}

func (r *repository) GetCampaignByCode(ctx context.Context, auth AuthContext, code string) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT campaign.get_campaign(p_code := $1)", code)
}

func (r *repository) UpsertCampaign(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT campaign.upsert_campaign($1::jsonb)", data)
}

func (r *repository) DeleteCampaign(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT campaign.delete_campaign($1)", id)
}

func (r *repository) RemoveHeroFromCampaign(ctx context.Context, auth AuthContext, heroID int64, campaignID int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT campaign.remove_hero_from_campaign($1, $2)", heroID, campaignID)
}

func (r *repository) GetCampaignSourceBookIDs(ctx context.Context, auth AuthContext, campaignID int64) ([]int64, error) {
	raw, err := r.callFunc(ctx, auth, "SELECT campaign.get_campaign_source_book_ids($1)", campaignID)
	if err != nil {
		return nil, err
	}
	var ids []int64
	if err := json.Unmarshal(raw, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}
