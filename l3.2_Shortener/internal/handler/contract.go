package handler

import (
	"context"

	"github.com/ProgrammistNik/WB-L3/l3.2/internal/dto"
	"github.com/ProgrammistNik/WB-L3/l3.2/internal/model"
)

type Service interface {
	Shorten(context.Context, dto.RequestURL) (*model.URL, error)
	GetOriginalURL(context.Context, string) (string, error)
	GetAnalytics(context.Context, string) ([]model.Click, error)
	TrackClick(ctx context.Context, shortURL, userAgent string) error
	GetAnalyticsGroupedByDay(context.Context, string) ([]model.AnalyticsResult, error)
	GetAnalyticsGroupedByMonth(context.Context, string) ([]model.AnalyticsResult, error)
	GetAnalyticsGroupedByUserAgent(context.Context, string) ([]model.AnalyticsResult, error)
}
