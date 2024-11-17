package repository

import (
	"context"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/entity"
)

type MovieRepository interface {
	Save(ctx context.Context, movie *entity.Movie) error
	FindByID(ctx context.Context, id string) (*entity.Movie, error)
}
