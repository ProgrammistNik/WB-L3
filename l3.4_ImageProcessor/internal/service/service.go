package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ProgrammistNik/WB-L3/l3.4_ImageProcessor/internal/kafka/producer"
	"github.com/ProgrammistNik/WB-L3/l3.4_ImageProcessor/internal/model"
	"github.com/disintegration/imaging"
)

type Service struct {
	storage     Storage
	producer    producer.ServiceProducer
	storagePath string
}

func New(st Storage, sp producer.ServiceProducer, storagePath string) *Service {
	os.MkdirAll(filepath.Join(storagePath, "originals"), os.ModePerm)
	os.MkdirAll(filepath.Join(storagePath, "processed"), os.ModePerm)
	os.MkdirAll(filepath.Join(storagePath, "thumbnails"), os.ModePerm)

	return &Service{
		storage:     st,
		producer:    sp,
		storagePath: storagePath,
	}
}

func (s *Service) UploadImage(ctx context.Context, file io.Reader, filename string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(filename + time.Now().String()))
	imageID := hex.EncodeToString(hash.Sum(nil))[:16]

	originalPath := filepath.Join(s.storagePath, "originals", imageID+filepath.Ext(filename))
	if err := s.storage.SaveFile(file, originalPath); err != nil {
		return "", fmt.Errorf("save file: %w", err)
	}

	imageData := &model.Image{
		ID:           imageID,
		OriginalPath: originalPath,
		Status:       model.StatusPending,
		CreatedAt:    time.Now(),
	}

	if err := s.storage.Create(ctx, imageData); err != nil {
		os.Remove(originalPath)
		return "", fmt.Errorf("create db record: %w", err)
	}

	message := map[string]string{
		"image_id": imageID,
		"path":     originalPath,
	}
	messageBytes, _ := json.Marshal(message)
	if err := s.producer.Send(ctx, []byte(imageID), messageBytes); err != nil {
		// Если Kafka не отправилась, удаляем запись и файл
		_ = s.storage.Delete(ctx, imageID)
		os.Remove(originalPath)
		return "", fmt.Errorf("send kafka message: %w", err)
	}

	return imageID, nil
}

func (s *Service) ProcessImage(ctx context.Context, imageID, imagePath string) error {
	imageData, err := s.storage.GetByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("get image: %w", err)
	}

	imageData.Status = model.StatusProcessing
	if err := s.storage.UpdateStatus(ctx, imageData); err != nil {
		return fmt.Errorf("update status processing: %w", err)
	}

	img, err := imaging.Open(imagePath)
	if err != nil {
		imageData.Status = model.StatusFailed
		_ = s.storage.UpdateStatus(ctx, imageData)
		return fmt.Errorf("open image: %w", err)
	}

	processedDir := filepath.Join(s.storagePath, "processed")
	thumbnailsDir := filepath.Join(s.storagePath, "thumbnails")

	resized := imaging.Resize(img, 1024, 0, imaging.Lanczos)
	resizedPath := filepath.Join(processedDir, imageID+"_resized.jpg")
	if err := imaging.Save(resized, resizedPath); err != nil {
		imageData.Status = model.StatusFailed
		_ = s.storage.UpdateStatus(ctx, imageData)
		return fmt.Errorf("save resized: %w", err)
	}

	thumb := imaging.Thumbnail(img, 200, 200, imaging.Lanczos)
	thumbPath := filepath.Join(thumbnailsDir, imageID+"_thumb.jpg")
	if err := imaging.Save(thumb, thumbPath); err != nil {
		imageData.Status = model.StatusFailed
		_ = s.storage.UpdateStatus(ctx, imageData)
		return fmt.Errorf("save thumbnail: %w", err)
	}

	wmPath := filepath.Join(s.storagePath, "watermark.png")
	final := resized
	if wm, err := imaging.Open(wmPath); err == nil {
		final = imaging.Overlay(resized, wm, image.Pt(20, 20), 0.4)
	}
	watermarkPath := filepath.Join(processedDir, imageID+"_watermarked.jpg")
	if err := imaging.Save(final, watermarkPath); err != nil {
		imageData.Status = model.StatusFailed
		_ = s.storage.UpdateStatus(ctx, imageData)
		return fmt.Errorf("save watermark: %w", err)
	}

	now := time.Now()
	imageData.Status = model.StatusCompleted
	imageData.ProcessedAt = sql.NullTime{Time: now, Valid: true}
	imageData.ResizedPath = sql.NullString{String: resizedPath, Valid: true}
	imageData.ThumbPath = sql.NullString{String: thumbPath, Valid: true}
	imageData.WatermarkPath = sql.NullString{String: watermarkPath, Valid: true}

	if err := s.storage.UpdateStatus(ctx, imageData); err != nil {
		return fmt.Errorf("update status completed: %w", err)
	}

	return nil
}

func (s *Service) GetImage(ctx context.Context, id string) (*model.Image, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *Service) DeleteImage(ctx context.Context, id string) error {
	imageData, err := s.storage.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.storage.Delete(ctx, id); err != nil {
		return err
	}

	filesToDelete := []string{
		imageData.OriginalPath,
		imageData.GetResizedPath(),
		imageData.GetThumbPath(),
		imageData.GetWatermarkPath(),
	}

	for _, f := range filesToDelete {
		if f != "" {
			_ = os.Remove(f)
		}
	}

	return nil
}
