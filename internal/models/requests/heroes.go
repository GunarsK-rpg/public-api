package requests

// GetHeroesQuery defines query params for GET /heroes.
type GetHeroesQuery struct {
	CampaignID *int64 `form:"campaign_id" json:"campaignId,omitempty"`
	Pager
}
