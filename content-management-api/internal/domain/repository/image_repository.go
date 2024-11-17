package repository

import (
	"context"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/entity"
)

type ImageRepository interface {
	Save(ctx context.Context, image *entity.Image) error
	FindByMovieIDAndType(ctx context.Context, movieID string, imageType entity.ImageType) (*entity.Image, error)
}
