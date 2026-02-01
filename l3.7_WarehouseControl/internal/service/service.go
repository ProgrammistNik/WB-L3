package service

import "github.com/ProgrammistNik/WB-L3/l3.7_WarehouseControl/internal/model"

type Service struct {
	storage Storage
}

func New(st Storage) *Service {
	return &Service{storage: st}
}

func (s *Service) CreateItem(username string, name string, qty int) error {
	return s.storage.CreateItem(name, qty, username)
}

func (s *Service) ListItems() ([]model.Item, error) {
	return s.storage.ListItems()
}

func (s *Service) UpdateItem(username string, id int, name string, qty int) error {
	return s.storage.UpdateItem(id, name, qty, username)
}

func (s *Service) DeleteItem(username string, id int) error {
	return s.storage.DeleteItem(id, username)
}

func (s *Service) GetHistory(id int) ([]model.ItemHistory, error) {
	return s.storage.GetHistory(id)
}
