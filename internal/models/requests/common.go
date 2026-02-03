package requests

// Pager defines common pagination parameters.
type Pager struct {
	Limit  *int `form:"limit" json:"limit,omitempty"`
	Offset int  `form:"offset" json:"offset"`
}
