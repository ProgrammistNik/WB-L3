package storage

import (
	"database/sql"
	"fmt"

	"github.com/ProgrammistNik/WB-L3/tree/main/l3.3_CommentTree/internal/model"
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

func (st *Storage) InsertComment(comment model.Comment) (int64, error) {
	query := `
		INSERT INTO comments (text, parent_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var id int64

	err := st.db.Master.QueryRow(query,
		comment.Text,
		comment.ParentID,
		comment.CreatedAt,
		comment.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("insert comment failed: %w", err)
	}

	return id, nil
}
func (st *Storage) GetTree(idComment string) ([]model.Comment, error) {
	var query string
	var rows *sql.Rows
	var err error

	zlog.Logger.Info().Str("idComment", idComment).Msg("Start GetTree")

	if idComment == "" {
		query = `
            WITH RECURSIVE comment_tree AS (
                SELECT * FROM comments WHERE parent_id IS NULL
                UNION ALL
                SELECT c.* FROM comments c
                INNER JOIN comment_tree ct ON c.parent_id = ct.id
            )
            SELECT * FROM comment_tree ORDER BY created_at;
        `
		rows, err = st.db.Master.Query(query)
	} else {
		query = `
            WITH RECURSIVE comment_tree AS (
                SELECT * FROM comments WHERE id = $1
                UNION ALL
                SELECT c.* FROM comments c
                INNER JOIN comment_tree ct ON c.parent_id = ct.id
            )
            SELECT * FROM comment_tree ORDER BY created_at;
        `
		rows, err = st.db.Master.Query(query, idComment)
	}

	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to execute recursive query")
		return nil, fmt.Errorf("get tree failed: %w", err)
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var c model.Comment
		err := rows.Scan(&c.ID, &c.ParentID, &c.Text, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("Failed to scan comment row")
			return nil, fmt.Errorf("scan comment failed: %w", err)
		}

		zlog.Logger.Debug().
			Int64("id", c.ID).
			Str("text", c.Text).
			Interface("parent_id", c.ParentID).
			Msg("Fetched comment")

		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		zlog.Logger.Error().Err(err).Msg("Rows iteration error")
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	zlog.Logger.Info().Int("total", len(comments)).Msg("Finished GetTree")
	return comments, nil
}

func (st *Storage) DeleteCommentByID(id string) error {
	query := `
		WITH RECURSIVE comment_tree AS (
			SELECT id FROM comments WHERE id = $1
			UNION ALL
			SELECT c.id FROM comments c
			INNER JOIN comment_tree ct ON c.parent_id = ct.id
		)
		DELETE FROM comments WHERE id IN (SELECT id FROM comment_tree);
	`

	_, err := st.db.Master.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete comment failed: %w", err)
	}

	return nil
}
func (st *Storage) SearchComments(query string, page, limit int) ([]model.Comment, error) {
	offset := (page - 1) * limit
	sqlQuery := `
		SELECT id, parent_id, text, created_at, updated_at
		FROM comments
		WHERE text ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := st.db.Master.Query(sqlQuery, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search comments failed: %w", err)
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var c model.Comment
		err := rows.Scan(&c.ID, &c.ParentID, &c.Text, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		comments = append(comments, c)
	}

	return comments, nil
}
