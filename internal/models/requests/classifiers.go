package requests

// GetClassifiersQuery defines query params for GET /classifiers.
type GetClassifiersQuery struct {
	CampaignID *int64 `form:"campaignId" json:"campaignId,omitempty"`
	HeroID     *int64 `form:"heroId" json:"heroId,omitempty"`
}
