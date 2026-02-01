package handler

import "github.com/ProgrammistNik/WB-L3/l3.7_WarehouseControl/internal/model"

type Service interface {
	CreateItem(string, string, int) error
	ListItems() ([]model.Item, error)
	UpdateItem(string, int, string, int) error
	DeleteItem(string, int) error
	GetHistory(int) ([]model.ItemHistory, error)
}
