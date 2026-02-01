package storage

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ProgrammistNik/WB-L3/l3.4_ImageProcessor/internal/model"
	"github.com/wb-go/wbf/dbpg"
)

type Storage struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (r *Storage) SaveFile(file io.Reader, path string) error {
	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func (r *Storage) Create(ctx context.Context, image *model.Image) error {
	query := `INSERT INTO images 
		(id, original_path, status, created_at) 
		VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query,
		image.ID, image.OriginalPath, image.Status, image.CreatedAt)
	return err
}

func (r *Storage) GetByID(ctx context.Context, imageID string) (*model.Image, error) {
	var image model.Image
	query := `SELECT id, original_path, resized_path, thumb_path, watermark_path, 
	                 status, created_at, processed_at
	          FROM images WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, imageID).
		Scan(&image.ID, &image.OriginalPath, &image.ResizedPath, &image.ThumbPath,
			&image.WatermarkPath, &image.Status, &image.CreatedAt, &image.ProcessedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}
	return &image, nil
}

func (r *Storage) UpdateStatus(ctx context.Context, image *model.Image) error {
	query := `UPDATE images 
		SET status = $1, processed_at = $2, resized_path = $3, 
		    thumb_path = $4, watermark_path = $5 
		WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query,
		image.Status, image.ProcessedAt, image.ResizedPath, image.ThumbPath, image.WatermarkPath, image.ID)
	return err
}

func (r *Storage) Delete(ctx context.Context, imageID string) error {
	query := `DELETE FROM images WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, imageID)
	return err
}
