package requests

import "fmt"

// GetHeroesQuery defines query params for GET /heroes.
type GetHeroesQuery struct {
	CampaignID *int64 `form:"campaign_id" json:"campaignId,omitempty"`
}

// Validate returns an error if any provided ID is non-positive.
func (q GetHeroesQuery) Validate() error {
	if q.CampaignID != nil && *q.CampaignID <= 0 {
		return fmt.Errorf("campaign_id must be a positive integer")
	}
	return nil
}
