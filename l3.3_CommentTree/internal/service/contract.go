package service

import (
	"github.com/ProgrammistNik/WB-L3/tree/main/l3.3_CommentTree/internal/model"
)

type Storage interface {
	InsertComment(comment model.Comment) (int64, error)
	GetTree(idComment string) ([]model.Comment, error)
	DeleteCommentByID(id string) error
	SearchComments(query string, page, limit int) ([]model.Comment, error)
}
