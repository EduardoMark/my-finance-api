package account

type AccountCreateRequest struct {
	Name    string  `json:"name" validate:"required"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
}
