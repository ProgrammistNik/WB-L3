package dto

type CommentRequest struct {
	Text     string `json:"text" validate:"required"`
	ParentID *int64 `json:"parent_id"`
}

type CommentResponse struct {
	ID       int64              `json:"id"`
	Text     string             `json:"text"`
	ParentID *int64             `json:"parent_id,omitempty"`
	Children []*CommentResponse `json:"children,omitempty"`
}
