package handler

import "github.com/ProgrammistNik/WB-L3/tree/main/l3.3_CommentTree/internal/handler/dto"

type Service interface {
	CreateComment(commentIn dto.CommentRequest) (dto.CommentResponse, error)
	GetAllComments(idComment string) ([]*dto.CommentResponse, error)
	DeleteComment(idComment string) error
	SearchComments(q string, page, limit int) ([]*dto.CommentResponse, error)
}
