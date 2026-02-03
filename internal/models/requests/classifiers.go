package requests

import "errors"

// ErrMissingActionLinkFilter is returned when neither object_id nor action_code is provided.
var ErrMissingActionLinkFilter = errors.New("either object_id or action_code is required")

// GetExpertisesQuery defines query params for GET /classifiers/expertises.
type GetExpertisesQuery struct {
	TypeCode *string `form:"type_code" json:"typeCode,omitempty"`
}

// GetSpecialtiesQuery defines query params for GET /classifiers/specialties.
type GetSpecialtiesQuery struct {
	PathCode *string `form:"path_code" json:"pathCode,omitempty"`
}

// GetSingerFormsQuery defines query params for GET /classifiers/singer-forms.
type GetSingerFormsQuery struct {
	BaseFormsOnly *bool `form:"base_forms_only" json:"baseFormsOnly,omitempty"`
}

// GetTalentsQuery defines query params for GET /classifiers/talents.
type GetTalentsQuery struct {
	PathCode         *string `form:"path_code" json:"pathCode,omitempty"`
	SpecialtyCode    *string `form:"specialty_code" json:"specialtyCode,omitempty"`
	AncestryCode     *string `form:"ancestry_code" json:"ancestryCode,omitempty"`
	RadiantOrderCode *string `form:"radiant_order_code" json:"radiantOrderCode,omitempty"`
	SurgeCode        *string `form:"surge_code" json:"surgeCode,omitempty"`
	IsKey            *bool   `form:"is_key" json:"isKey,omitempty"`
}

// GetActionsQuery defines query params for GET /classifiers/actions.
type GetActionsQuery struct {
	ActionTypeCode     *string `form:"action_type_code" json:"actionTypeCode,omitempty"`
	ActivationTypeCode *string `form:"activation_type_code" json:"activationTypeCode,omitempty"`
	DamageTypeCode     *string `form:"damage_type_code" json:"damageTypeCode,omitempty"`
}

// GetActionLinksQuery defines query params for GET /classifiers/action-links.
type GetActionLinksQuery struct {
	ObjectID   *int64  `form:"object_id" json:"objectId,omitempty"`
	ActionCode *string `form:"action_code" json:"actionCode,omitempty"`
}

// Validate checks that at least one filter is provided.
func (q GetActionLinksQuery) Validate() error {
	if q.ObjectID == nil && q.ActionCode == nil {
		return ErrMissingActionLinkFilter
	}
	return nil
}

// GetEquipmentsQuery defines query params for GET /classifiers/equipments.
type GetEquipmentsQuery struct {
	TypeCode       *string `form:"type_code" json:"typeCode,omitempty"`
	DamageTypeCode *string `form:"damage_type_code" json:"damageTypeCode,omitempty"`
}
