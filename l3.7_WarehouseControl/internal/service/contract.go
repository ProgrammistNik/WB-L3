package service

import "github.com/ProgrammistNik/WB-L3/l3.7_WarehouseControl/internal/model"

type Storage interface {
	CreateItem(string, int, string) error
	ListItems() ([]model.Item, error)
	GetHistory(id int) ([]model.ItemHistory, error)
	DeleteItem(id int, username string) error
	UpdateItem(id int, name string, qty int, username string) error
}
