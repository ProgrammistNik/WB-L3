package service

import (
	"errors"
	"testing"

	"github.com/ProgrammistNik/WB-L3/tree/main/l3.3_CommentTree/internal/handler/dto"
	"github.com/ProgrammistNik/WB-L3/tree/main/l3.3_CommentTree/internal/model"
	"github.com/stretchr/testify/assert"
)

type mockStorage struct {
	InsertCommentFunc     func(comment model.Comment) (int64, error)
	GetTreeFunc           func(idComment string) ([]model.Comment, error)
	DeleteCommentByIDFunc func(id string) error
	SearchCommentsFunc    func(query string, page, limit int) ([]model.Comment, error)
}

func (m *mockStorage) InsertComment(comment model.Comment) (int64, error) {
	return m.InsertCommentFunc(comment)
}

func (m *mockStorage) GetTree(idComment string) ([]model.Comment, error) {
	return m.GetTreeFunc(idComment)
}

func (m *mockStorage) DeleteCommentByID(id string) error {
	return m.DeleteCommentByIDFunc(id)
}

func (m *mockStorage) SearchComments(query string, page, limit int) ([]model.Comment, error) {
	return m.SearchCommentsFunc(query, page, limit)
}

func TestService_CreateComment(t *testing.T) {
	mock := &mockStorage{
		InsertCommentFunc: func(comment model.Comment) (int64, error) {
			return 42, nil
		},
	}
	svc := New(mock)

	req := dto.CommentRequest{
		Text:     "Hello",
		ParentID: nil,
	}

	res, err := svc.CreateComment(req)

	assert.NoError(t, err)
	assert.Equal(t, int64(42), res.ID)
	assert.Equal(t, "Hello", res.Text)
}

func TestService_GetAllComments(t *testing.T) {
	mock := &mockStorage{
		GetTreeFunc: func(id string) ([]model.Comment, error) {
			return []model.Comment{
				{ID: 1, Text: "Root"},
				{ID: 2, Text: "Child", ParentID: ptrInt64(1)},
			}, nil
		},
	}
	svc := New(mock)

	res, err := svc.GetAllComments("")

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "Root", res[0].Text)
	assert.Len(t, res[0].Children, 1)
	assert.Equal(t, "Child", res[0].Children[0].Text)
}

func ptrInt64(i int64) *int64 {
	return &i
}

func TestService_DeleteComment(t *testing.T) {
	mock := &mockStorage{
		DeleteCommentByIDFunc: func(id string) error {
			if id == "42" {
				return nil
			}
			return errors.New("not found")
		},
	}
	svc := New(mock)

	err := svc.DeleteComment("42")
	assert.NoError(t, err)

	err = svc.DeleteComment("999")
	assert.Error(t, err)
}

func TestService_SearchComments(t *testing.T) {
	mock := &mockStorage{
		GetTreeFunc: func(id string) ([]model.Comment, error) {
			return []model.Comment{
				{ID: 1, Text: "Root"},
				{ID: 2, Text: "Hello world", ParentID: ptrInt64(1)},
			}, nil
		},
	}
	svc := New(mock)

	res, err := svc.SearchComments("hello", 1, 10)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "Root", res[0].Text)
	assert.Len(t, res[0].Children, 1)
	assert.Equal(t, "Hello world", res[0].Children[0].Text)
}
