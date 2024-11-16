package database

import (
	"context"
	"fmt"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGVideoRepository struct {
	pool *pgxpool.Pool
}

func (pg *PGVideoRepository) Save(ctx context.Context, video *entity.Video) error {
	query := `
		INSERT INTO videos (
			id, title, description, file_path, storage_location, content_type, 
			storage_provider, original_file_name, created_at, updated_at, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			file_path = EXCLUDED.file_path,
			storage_location = EXCLUDED.storage_location,
			content_type = EXCLUDED.content_type,
			storage_provider = EXCLUDED.storage_provider,
			original_file_name = EXCLUDED.original_file_name,
			updated_at = NOW(),
			status = EXCLUDED.status
	`

	_, err := pg.pool.Exec(ctx, query,
		video.ID,
		video.Title,
		video.Description,
		video.FilePath,
		video.StorageLocation,
		video.ContentType,
		video.StorageProvider,
		video.OriginalFileName,
		video.CreatedAt,
		video.UpdatedAt,
		video.Status,
	)

	if err != nil {
		return fmt.Errorf("failed to save video: %w", err)
	}

	return nil
}

func NewPGVideoRepository(pool *pgxpool.Pool) *PGVideoRepository {
	return &PGVideoRepository{
		pool,
	}
}
