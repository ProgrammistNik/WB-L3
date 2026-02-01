package model

import "github.com/ProgrammistNik/WB-L3/tree/main/l3.3_CommentTree/internal/handler/dto"

func CastModel(c Comment) *dto.CommentResponse {
	return &dto.CommentResponse{
		ID:       c.ID,
		Text:     c.Text,
		ParentID: c.ParentID,
		Children: []*dto.CommentResponse{},
	}
}
