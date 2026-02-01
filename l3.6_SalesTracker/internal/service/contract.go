package service

import (
	"github.com/ProgrammistNik/WB-L3/l3.6_SalesTracker/internal/model"
	"github.com/wb-go/wbf/ginext"
)

type Storage interface {
	SaveItem(*ginext.Context, model.Item) (model.Item, error)
	AnalyticsCalculate(*ginext.Context, model.ItemsFilter) (model.AnalyticsResponse, error)
	GetItems(*ginext.Context, model.ItemsFilter) ([]model.Item, error)
	UpdateItem(*ginext.Context, model.Item) (model.Item, error)
	DeleteItem(*ginext.Context, int64) error
}
