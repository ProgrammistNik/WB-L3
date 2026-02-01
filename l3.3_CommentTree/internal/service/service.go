package service

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/ProgrammistNik/WB-L3/tree/main/l3.3_CommentTree/internal/handler/dto"
	"github.com/ProgrammistNik/WB-L3/tree/main/l3.3_CommentTree/internal/model"
	"github.com/wb-go/wbf/zlog"
)

type Service struct {
	storage Storage
}

func New(st Storage) *Service {
	return &Service{
		storage: st,
	}
}

func (s *Service) CreateComment(commentIncome dto.CommentRequest) (dto.CommentResponse, error) {
	comment := model.Comment{
		Text:      commentIncome.Text,
		ParentID:  commentIncome.ParentID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := s.storage.InsertComment(comment)
	if err != nil {
		return dto.CommentResponse{}, err
	}

	comment.ID = id
	return *model.CastModel(comment), nil
}

func (s *Service) GetAllComments(idComment string) ([]*dto.CommentResponse, error) {
	comments, err := s.storage.GetTree(idComment)
	if err != nil {
		return nil, err
	}
	tree := buildTree(comments)

	// DEBUG
	if debugJson, err := json.MarshalIndent(tree, "", "  "); err == nil {
		zlog.Logger.Info().Msg("Full comment tree:\n" + string(debugJson))
	}

	return tree, nil
}

func (s *Service) SearchComments(q string, page, limit int) ([]*dto.CommentResponse, error) {
	allComments, err := s.storage.GetTree("")
	if err != nil {
		return nil, err
	}

	tree := buildTree(allComments)
	var result []*dto.CommentResponse
	for _, root := range tree {
		if found := findCommentsByKeyword(root, q); found != nil {
			result = append(result, found)
		}
	}
	return result, nil
}

func findCommentsByKeyword(c *dto.CommentResponse, keyword string) *dto.CommentResponse {
	var matchingChildren []*dto.CommentResponse

	for _, child := range c.Children {
		if match := findCommentsByKeyword(child, keyword); match != nil {
			matchingChildren = append(matchingChildren, match)
		}
	}

	if strings.Contains(strings.ToLower(c.Text), strings.ToLower(keyword)) || len(matchingChildren) > 0 {
		return &dto.CommentResponse{
			ID:       c.ID,
			Text:     c.Text,
			Children: matchingChildren,
		}
	}

	return nil
}

func buildTree(comments []model.Comment) []*dto.CommentResponse {
	idToComment := make(map[int64]*dto.CommentResponse)

	for _, c := range comments {
		copied := model.CastModel(c)
		copied.Children = []*dto.CommentResponse{}
		idToComment[c.ID] = copied
	}

	var roots []*dto.CommentResponse
	for _, c := range comments {
		current := idToComment[c.ID]

		if c.ParentID != nil {
			parent, ok := idToComment[*c.ParentID]
			if ok {
				parent.Children = append(parent.Children, current)
			}
		} else {
			roots = append(roots, current)
		}
	}

	return roots
}

func (s *Service) DeleteComment(id string) error {
	err := s.storage.DeleteCommentByID(id)
	if err != nil {
		return err
	}

	return nil
}
