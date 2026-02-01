package service

import (
	"context"
	"io"

	"github.com/ProgrammistNik/WB-L3/l3.4_ImageProcessor/internal/model"
)

type Storage interface {
	SaveFile(io.Reader, string) error
	Create(context.Context, *model.Image) error
	GetByID(context.Context, string) (*model.Image, error)
	UpdateStatus(context.Context, *model.Image) error
	Delete(context.Context, string) error
}
