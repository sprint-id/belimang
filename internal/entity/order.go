package entity

type (
	Order struct {
		ID         string `json:"id"`
		EstimateID string `json:"estimate_id"`

		UserID string `json:"user_id"`
	}
)
