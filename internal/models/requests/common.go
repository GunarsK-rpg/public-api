package requests

// Pager defines common pagination parameters.
type Pager struct {
	Limit  *int `form:"limit" json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
	Offset int  `form:"offset" json:"offset" validate:"min=0"`
}
