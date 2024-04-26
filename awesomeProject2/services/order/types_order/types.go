package types_order

type OrderPayload struct {
	UserId  string `json:"userID" `
	Total   int    `json:"total" validate:"required"`
	Status  string `json:"status" validate:"required"`
	Address string `json:"address" validate:"required"`
}
