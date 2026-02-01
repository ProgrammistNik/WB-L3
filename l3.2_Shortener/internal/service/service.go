package service

import (
	"context"
	"encoding/base64"
	"errors"
	"net/url"
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.2/internal/dto"
	"github.com/ProgrammistNik/WB-L3/l3.2/internal/model"
	"github.com/wb-go/wbf/zlog"
	"golang.org/x/exp/rand"
)

type Service struct {
	storage Storage
}

func New(st Storage) *Service {
	return &Service{
		storage: st,
	}
}

func (s *Service) Shorten(ctx context.Context, urlRequest dto.RequestURL) (*model.URL, error) {
	if !isValidURL(urlRequest.URL) {
		err := errors.New("invalid url format")
		zlog.Logger.Error().Err(err).Str("url", urlRequest.URL).Msg("Invalid URL")
		return nil, err
	}

	code, err := s.generateUniqueShortCode(ctx)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to generate unique short code")
		return nil, err
	}

	link := &model.URL{
		OriginalURL: urlRequest.URL,
		ShortURL:    code,
		CreateAt:    time.Now(),
	}

	if err := s.storage.SaveLink(ctx, link); err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to save link to storage")
		return nil, err
	}

	return link, nil
}

func isValidURL(raw string) bool {
	u, err := url.ParseRequestURI(raw)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (s *Service) generateUniqueShortCode(ctx context.Context) (string, error) {
	const maxAttempts = 5

	for i := 0; i < maxAttempts; i++ {
		code := generateShortURL()
		zlog.Logger.Info().Str("code", code).Msg("Generated short code")

		exists, err := s.storage.ExistsByShortCode(ctx, code)
		if err != nil {
			return "", err
		}
		zlog.Logger.Info().Bool("exists", exists).Msg("Exists check result")

		if !exists {
			return code, nil
		}
	}

	return "", errors.New("failed to generate unique short code after several attempts")
}

func generateShortURL() string {
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		zlog.Logger.Error().Err(err).Msg("rand.Read error")
		return ""
	}
	code := base64.URLEncoding.EncodeToString(b)[:7]
	zlog.Logger.Info().Str("short_url", code).Msg("Generated short URL")
	return code
}

func (s *Service) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	url, err := s.storage.GetOriginalURL(ctx, shortCode)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *Service) GetAnalytics(ctx context.Context, shortURL string) ([]model.Click, error) {
	linkID, err := s.storage.GetLinkIDByShortURL(ctx, shortURL)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("short_url", shortURL).
			Msg("Failed to get link ID by short URL")
		return nil, err
	}

	clicks, err := s.storage.GetClicksByLinkID(ctx, linkID)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("link_id", linkID).
			Msg("Failed to get clicks for link ID")
		return nil, err
	}

	return clicks, nil
}

func (s *Service) TrackClick(ctx context.Context, shortURL, userAgent string) error {
	linkID, err := s.storage.GetLinkIDByShortURL(ctx, shortURL)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Str("short_url", shortURL).
			Msg("Failed to get link ID to track click")
		return err
	}

	err = s.storage.InsertClick(ctx, linkID, userAgent)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("link_id", linkID).
			Str("short_url", shortURL).
			Msg("Failed to insert click")
		return err
	}

	return nil
}

func (s *Service) GetAnalyticsGroupedByDay(ctx context.Context, shortURL string) ([]model.AnalyticsResult, error) {
	linkID, err := s.storage.GetLinkIDByShortURL(ctx, shortURL)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("short_url", shortURL).Msg("GetLinkIDByShortURL failed (group=day)")
		return nil, err
	}

	results, err := s.storage.GetClicksGroupedByDay(ctx, linkID)
	if err != nil {
		zlog.Logger.Error().Err(err).Int("link_id", linkID).Msg("GetClicksGroupedByDay failed")
		return nil, err
	}

	return results, nil
}

func (s *Service) GetAnalyticsGroupedByMonth(ctx context.Context, shortURL string) ([]model.AnalyticsResult, error) {
	linkID, err := s.storage.GetLinkIDByShortURL(ctx, shortURL)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("short_url", shortURL).Msg("GetLinkIDByShortURL failed (group=month)")
		return nil, err
	}

	results, err := s.storage.GetClicksGroupedByMonth(ctx, linkID)
	if err != nil {
		zlog.Logger.Error().Err(err).Int("link_id", linkID).Msg("GetClicksGroupedByMonth failed")
		return nil, err
	}

	return results, nil
}

func (s *Service) GetAnalyticsGroupedByUserAgent(ctx context.Context, shortURL string) ([]model.AnalyticsResult, error) {
	linkID, err := s.storage.GetLinkIDByShortURL(ctx, shortURL)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("short_url", shortURL).Msg("GetLinkIDByShortURL failed (group=ua)")
		return nil, err
	}

	results, err := s.storage.GetClicksGroupedByUserAgent(ctx, linkID)
	if err != nil {
		zlog.Logger.Error().Err(err).Int("link_id", linkID).Msg("GetClicksGroupedByUserAgent failed")
		return nil, err
	}

	return results, nil
}
