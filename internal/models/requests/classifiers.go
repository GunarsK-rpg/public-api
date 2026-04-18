package requests

import "fmt"

// GetClassifiersQuery defines query params for GET /classifiers.
// campaignId and heroId may be combined (character-sheet scope).
// sourceBookId is the Library scope and cannot be combined with the others.
type GetClassifiersQuery struct {
	CampaignID   *int64 `form:"campaignId" json:"campaignId,omitempty"`
	HeroID       *int64 `form:"heroId" json:"heroId,omitempty"`
	SourceBookID *int64 `form:"sourceBookId" json:"sourceBookId,omitempty"`
}

// Validate returns an error if any provided ID is non-positive or sourceBookId
// is combined with campaignId / heroId.
func (q GetClassifiersQuery) Validate() error {
	if q.CampaignID != nil && *q.CampaignID <= 0 {
		return fmt.Errorf("campaignId must be a positive integer")
	}
	if q.HeroID != nil && *q.HeroID <= 0 {
		return fmt.Errorf("heroId must be a positive integer")
	}
	if q.SourceBookID != nil && *q.SourceBookID <= 0 {
		return fmt.Errorf("sourceBookId must be a positive integer")
	}
	if q.SourceBookID != nil && (q.CampaignID != nil || q.HeroID != nil) {
		return fmt.Errorf("sourceBookId cannot be combined with campaignId or heroId")
	}
	return nil
}
