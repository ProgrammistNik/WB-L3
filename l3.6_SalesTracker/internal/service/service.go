package service

import (
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.6_SalesTracker/internal/model"
	"github.com/wb-go/wbf/ginext"
)

type Service struct {
	storage Storage
}

func New(st Storage) *Service {
	return &Service{storage: st}
}

func (s *Service) CreateItem(c *ginext.Context, item model.Item) (model.Item, error) {
	item.CreatedAt = time.Now()

	data, err := s.storage.SaveItem(c, item)
	if err != nil {
		return model.Item{}, err
	}

	return data, nil
}

func (s *Service) GetAnalytics(c *ginext.Context, filter model.ItemsFilter) (model.AnalyticsResponse, error) {
	return s.storage.AnalyticsCalculate(c, filter)
}

func (s *Service) GetItems(c *ginext.Context, filter model.ItemsFilter) ([]model.Item, error) {
	return s.storage.GetItems(c, filter)
}

func (s *Service) UpdateItem(c *ginext.Context, item model.Item) (model.Item, error) {
	return s.storage.UpdateItem(c, item)
}

func (s *Service) DeleteItem(c *ginext.Context, id int64) error {
	return s.storage.DeleteItem(c, id)
}
