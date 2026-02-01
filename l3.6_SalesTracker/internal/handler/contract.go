package handler

import (
	"github.com/ProgrammistNik/WB-L3/l3.6_SalesTracker/internal/model"
	"github.com/wb-go/wbf/ginext"
)

type Service interface {
	CreateItem(*ginext.Context, model.Item) (model.Item, error)
	GetAnalytics(*ginext.Context, model.ItemsFilter) (model.AnalyticsResponse, error)
	GetItems(*ginext.Context, model.ItemsFilter) ([]model.Item, error)
	UpdateItem(*ginext.Context, model.Item) (model.Item, error)
	DeleteItem(*ginext.Context, int64) error
}
