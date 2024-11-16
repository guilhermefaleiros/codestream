package repository

import (
	"context"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/entity"
)

type VideoRepository interface {
	Save(ctx context.Context, video *entity.Video) error
}
