package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ProgrammistNik/WB-L3/l3.2/internal/model"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

type Storage struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (st *Storage) SaveLink(ctx context.Context, url *model.URL) error {
	query := `
		INSERT INTO links (original_url, short_url, created_at)
		VALUES ($1, $2, $3)
	`
	_, err := st.db.ExecContext(ctx, query, url.OriginalURL, url.ShortURL, url.CreateAt)
	return err
}

func (st *Storage) ExistsByShortCode(ctx context.Context, shortCode string) (bool, error) {
	query := `SELECT 1 FROM links WHERE short_url = $1 LIMIT 1`

	row := st.db.Master.QueryRowContext(ctx, query, shortCode)

	var exists int
	err := row.Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (st *Storage) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	query := `SELECT original_url FROM links WHERE short_url = $1`

	var originalURL string
	err := st.db.Master.QueryRowContext(ctx, query, shortCode).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("not found")
		}
		return "", err
	}

	return originalURL, nil
}

func (st *Storage) GetLinkIDByShortURL(ctx context.Context, shortURL string) (int, error) {
	query := `SELECT id FROM links WHERE short_url = $1 LIMIT 1`

	var id int
	err := st.db.Master.QueryRowContext(ctx, query, shortURL).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("short url not found")
		}
		return 0, err
	}

	return id, nil
}

func (st *Storage) GetClicksByLinkID(ctx context.Context, linkID int) ([]model.Click, error) {
	query := `
		SELECT id, link_id, user_agent, created_at
		FROM clicks
		WHERE link_id = $1
		ORDER BY created_at DESC
	`

	rows, err := st.db.Master.QueryContext(ctx, query, linkID)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("link_id", linkID).
			Msg("Failed to query clicks from database")
		return nil, err
	}
	defer rows.Close()

	var clicks []model.Click
	for rows.Next() {
		var c model.Click
		err := rows.Scan(&c.ID, &c.LinkID, &c.UserAgent, &c.CreateAt)
		if err != nil {
			zlog.Logger.Error().
				Err(err).
				Msg("Failed to scan click row")
			return nil, err
		}
		clicks = append(clicks, c)
	}

	if err = rows.Err(); err != nil {
		zlog.Logger.Error().
			Err(err).
			Msg("Row iteration error")
		return nil, err
	}

	return clicks, nil
}

func (st *Storage) InsertClick(ctx context.Context, linkID int, userAgent string) error {
	query := `
		INSERT INTO clicks (link_id, user_agent, created_at)
		VALUES ($1, $2, NOW())
	`
	_, err := st.db.ExecContext(ctx, query, linkID, userAgent)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("link_id", linkID).
			Str("user_agent", userAgent).
			Msg("Failed to insert click")
	}
	return err
}

func (st *Storage) GetClicksGroupedByDay(ctx context.Context, linkID int) ([]model.AnalyticsResult, error) {
	query := `
		SELECT TO_CHAR(created_at, 'YYYY-MM-DD') AS group_by, COUNT(*) AS count
		FROM clicks
		WHERE link_id = $1
		GROUP BY group_by
		ORDER BY group_by
	`
	return st.queryAnalytics(ctx, query, linkID)
}

func (st *Storage) GetClicksGroupedByMonth(ctx context.Context, linkID int) ([]model.AnalyticsResult, error) {
	query := `
		SELECT TO_CHAR(created_at, 'YYYY-MM') AS group_by, COUNT(*) AS count
		FROM clicks
		WHERE link_id = $1
		GROUP BY group_by
		ORDER BY group_by
	`
	return st.queryAnalytics(ctx, query, linkID)
}

func (st *Storage) GetClicksGroupedByUserAgent(ctx context.Context, linkID int) ([]model.AnalyticsResult, error) {
	query := `
		SELECT user_agent AS group_by, COUNT(*) AS count
		FROM clicks
		WHERE link_id = $1
		GROUP BY group_by
		ORDER BY count DESC
	`
	return st.queryAnalytics(ctx, query, linkID)
}

func (st *Storage) queryAnalytics(ctx context.Context, query string, linkID int) ([]model.AnalyticsResult, error) {
	rows, err := st.db.Master.QueryContext(ctx, query, linkID)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to execute analytics query")
		return nil, err
	}
	defer rows.Close()

	var results []model.AnalyticsResult
	for rows.Next() {
		var res model.AnalyticsResult
		if err := rows.Scan(&res.Group, &res.Count); err != nil {
			zlog.Logger.Error().Err(err).Msg("Failed to scan analytics result row")
			return nil, err
		}
		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		zlog.Logger.Error().Err(err).Msg("Row iteration error in analytics")
		return nil, err
	}

	return results, nil
}
