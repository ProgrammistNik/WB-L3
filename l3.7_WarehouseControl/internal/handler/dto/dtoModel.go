package dto

type CreateItemRequest struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type UpdateItemRequest struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
