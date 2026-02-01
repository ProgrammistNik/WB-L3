package dto

type ItemRequest struct {
	Type     string  `json:"type" binding:"required"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount" binding:"required,gte=0"`
	Date     string  `json:"date" binding:"required"`
}

type AnalyticsRequest struct {
	From     *string `form:"from"`
	To       *string `form:"to"`
	Category *string `form:"category"`
}

type GetItemsRequest struct {
	From     *string `form:"from"`
	To       *string `form:"to"`
	Category *string `form:"category"`
	Limit    *int    `form:"limit"`
	Offset   *int    `form:"offset"`
}

type UpdateItemRequest struct {
	Type     *string  `json:"type"`
	Category *string  `json:"category"`
	Amount   *float64 `json:"amount"`
	Date     *string  `json:"date"`
}
