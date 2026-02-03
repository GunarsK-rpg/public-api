package repository

import (
	"context"
	"encoding/json"
)

// CampaignRepository defines methods for campaign data access.
type CampaignRepository interface {
	GetCampaigns(ctx context.Context, auth AuthContext) (json.RawMessage, error)
	GetCampaign(ctx context.Context, auth AuthContext, id int64) (json.RawMessage, error)
	UpsertCampaign(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error)
	DeleteCampaign(ctx context.Context, auth AuthContext, id int64) (bool, error)
}

func (r *repository) GetCampaigns(ctx context.Context, auth AuthContext) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT campaigns.get_campaigns()")
}

func (r *repository) GetCampaign(ctx context.Context, auth AuthContext, id int64) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT campaigns.get_campaign($1)", id)
}

func (r *repository) UpsertCampaign(ctx context.Context, auth AuthContext, data json.RawMessage) (json.RawMessage, error) {
	return r.callFunc(ctx, auth, "SELECT campaigns.upsert_campaign($1::jsonb)", data)
}

func (r *repository) DeleteCampaign(ctx context.Context, auth AuthContext, id int64) (bool, error) {
	return r.execFunc(ctx, auth, "SELECT campaigns.delete_campaign($1)", id)
}
