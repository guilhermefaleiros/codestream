package database

import (
	"context"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGImageRepository struct {
	pool *pgxpool.Pool
}

func (pg *PGImageRepository) Save(ctx context.Context, image *entity.Image) error {
	query := `
		INSERT INTO images (
			id, file_path, storage_location, content_type, storage_provider, original_file_name, movie_id, type, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
		ON CONFLICT (id) DO UPDATE SET
			file_path = EXCLUDED.file_path,
			storage_location = EXCLUDED.storage_location,
			content_type = EXCLUDED.content_type,
			storage_provider = EXCLUDED.storage_provider,
			original_file_name = EXCLUDED.original_file_name,
			movie_id = EXCLUDED.movie_id,
			type = EXCLUDED.type,
			updated_at = NOW()
	`

	_, err := pg.pool.Exec(ctx, query,
		image.ID,
		image.FilePath,
		image.StorageLocation,
		image.ContentType,
		image.StorageProvider,
		image.OriginalFileName,
		image.MovieID,
		image.Type,
		image.CreatedAt,
		image.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (pg *PGImageRepository) FindByMovieIDAndType(ctx context.Context, movieID string, imageType entity.ImageType) (*entity.Image, error) {
	query := `
		SELECT id, file_path, storage_location, content_type, storage_provider, original_file_name, movie_id, type, created_at, updated_at
		FROM images
		WHERE movie_id = $1 AND type = $2
	`

	row := pg.pool.QueryRow(ctx, query, movieID, imageType)

	var image entity.Image
	err := row.Scan(
		&image.ID,
		&image.FilePath,
		&image.StorageLocation,
		&image.ContentType,
		&image.StorageProvider,
		&image.OriginalFileName,
		&image.MovieID,
		&image.Type,
		&image.CreatedAt,
		&image.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &image, nil
}

func NewPGImageRepository(pool *pgxpool.Pool) *PGImageRepository {
	return &PGImageRepository{
		pool: pool,
	}
}
