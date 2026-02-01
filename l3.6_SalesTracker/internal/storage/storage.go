package storage

import (
	"strconv"
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.6_SalesTracker/internal/model"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
)

type Storage struct{ db *dbpg.DB }

func New(db *dbpg.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (st *Storage) SaveItem(c *ginext.Context, item model.Item) (model.Item, error) {
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO items (type, category, amount, date, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := st.db.Master.QueryRow(query,
		item.Type,
		item.Category,
		item.Amount,
		item.Date,
		item.CreatedAt,
	).Scan(&item.ID, &item.CreatedAt)
	if err != nil {
		return model.Item{}, err
	}

	return item, nil
}

func (st *Storage) AnalyticsCalculate(c *ginext.Context, filter model.ItemsFilter) (model.AnalyticsResponse, error) {
	query := `
        SELECT 
            COALESCE(SUM(amount), 0) as total_sum,
            COALESCE(AVG(amount), 0) as average,
            COUNT(amount) as count,
            COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount), 0) AS median,
            COALESCE(PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount), 0) AS percentile_90
        FROM items 
        WHERE 1=1
    `

	args := []any{}
	argCounter := 1

	if filter.From != nil {
		query += " AND date >= $" + strconv.Itoa(argCounter)
		args = append(args, *filter.From)
		argCounter++
	}
	if filter.To != nil {
		query += " AND date <= $" + strconv.Itoa(argCounter)
		args = append(args, *filter.To)
		argCounter++
	}
	if filter.Category != nil && *filter.Category != "" {
		query += " AND category = $" + strconv.Itoa(argCounter)
		args = append(args, *filter.Category)
		argCounter++
	}
	_ = argCounter

	var response model.AnalyticsResponse
	err := st.db.Master.QueryRow(query, args...).Scan(
		&response.Sum,
		&response.Avg,
		&response.Count,
		&response.Median,
		&response.P90,
	)
	if err != nil {
		return model.AnalyticsResponse{}, err
	}

	return response, nil
}

func (st *Storage) GetItems(c *ginext.Context, filter model.ItemsFilter) ([]model.Item, error) {
	query := "SELECT id, type, category, amount, date, created_at FROM items WHERE 1=1"
	args := []any{}
	argCounter := 1

	if filter.From != nil {
		query += " AND date >= $" + strconv.Itoa(argCounter)
		args = append(args, *filter.From)
		argCounter++
	}
	if filter.To != nil {
		query += " AND date <= $" + strconv.Itoa(argCounter)
		args = append(args, *filter.To)
		argCounter++
	}
	if filter.Category != nil && *filter.Category != "" {
		query += " AND category = $" + strconv.Itoa(argCounter)
		args = append(args, *filter.Category)
		argCounter++
	}

	if filter.Limit != nil {
		query += " LIMIT $" + strconv.Itoa(argCounter)
		args = append(args, *filter.Limit)
		argCounter++
	}
	if filter.Offset != nil {
		query += " OFFSET $" + strconv.Itoa(argCounter)
		args = append(args, *filter.Offset)
		argCounter++
	}

	_ = argCounter

	rows, err := st.db.Master.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var it model.Item
		if err := rows.Scan(&it.ID, &it.Type, &it.Category, &it.Amount, &it.Date, &it.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, it)
	}

	return items, nil
}

func (st *Storage) UpdateItem(c *ginext.Context, item model.Item) (model.Item, error) {
	query := `
		UPDATE items
		SET type = COALESCE(NULLIF($1, ''), type),
		    category = COALESCE(NULLIF($2, ''), category),
		    amount = COALESCE($3, amount),
		    date = COALESCE($4, date)
		WHERE id = $5
		RETURNING id, type, category, amount, date, created_at
	`

	var updated model.Item
	err := st.db.Master.QueryRow(query,
		item.Type,
		item.Category,
		item.Amount,
		item.Date,
		item.ID,
	).Scan(
		&updated.ID,
		&updated.Type,
		&updated.Category,
		&updated.Amount,
		&updated.Date,
		&updated.CreatedAt,
	)
	if err != nil {
		return model.Item{}, err
	}

	return updated, nil
}

func (st *Storage) DeleteItem(c *ginext.Context, id int64) error {
	query := `DELETE FROM items WHERE id = $1`
	_, err := st.db.Master.Exec(query, id)
	return err
}
