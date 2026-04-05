package requests

import "fmt"

// GetClassifiersQuery defines query params for GET /classifiers.
type GetClassifiersQuery struct {
	CampaignID *int64 `form:"campaignId" json:"campaignId,omitempty"`
	HeroID     *int64 `form:"heroId" json:"heroId,omitempty"`
}

// Validate returns an error if any provided ID is non-positive.
func (q GetClassifiersQuery) Validate() error {
	if q.CampaignID != nil && *q.CampaignID <= 0 {
		return fmt.Errorf("campaignId must be a positive integer")
	}
	if q.HeroID != nil && *q.HeroID <= 0 {
		return fmt.Errorf("heroId must be a positive integer")
	}
	return nil
}
