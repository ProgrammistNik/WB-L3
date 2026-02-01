package service

import (
	"context"

	"github.com/ProgrammistNik/WB-L3/l3.2/internal/model"
)

type Storage interface {
	SaveLink(context.Context, *model.URL) error
	ExistsByShortCode(context.Context, string) (bool, error)
	GetOriginalURL(context.Context, string) (string, error)
	GetLinkIDByShortURL(ctx context.Context, shortURL string) (int, error)
	GetClicksByLinkID(ctx context.Context, linkID int) ([]model.Click, error)
	InsertClick(ctx context.Context, linkID int, userAgent string) error
	GetClicksGroupedByDay(ctx context.Context, linkID int) ([]model.AnalyticsResult, error)
	GetClicksGroupedByMonth(ctx context.Context, linkID int) ([]model.AnalyticsResult, error)
	GetClicksGroupedByUserAgent(ctx context.Context, linkID int) ([]model.AnalyticsResult, error)
}
