package database

import (
	"context"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGMovieRepository struct {
	pool *pgxpool.Pool
}

func (pg *PGMovieRepository) Save(ctx context.Context, movie *entity.Movie) error {
	query := `
		INSERT INTO movies (
			id, title, description, launch_year, genre, duration, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		ON CONFLICT (id) DO UPDATE SET
				
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			launch_year = EXCLUDED.launch_year,
			genre = EXCLUDED.genre,
			duration = EXCLUDED.duration,
			status = EXCLUDED.status,
			updated_at = NOW()
		                    
	`

	_, err := pg.pool.Exec(ctx, query,
		movie.ID,
		movie.Title,
		movie.Description,
		movie.LaunchYear,
		movie.Genre,
		movie.Duration,
		movie.Status,
		movie.CreatedAt,
		movie.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (pg *PGMovieRepository) FindByID(ctx context.Context, movieID string) (*entity.Movie, error) {
	query := `
		SELECT id, title, description, launch_year, genre, duration, status, created_at, updated_at
		FROM movies
		WHERE id = $1
	`

	var movie entity.Movie
	err := pg.pool.QueryRow(ctx, query, movieID).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.LaunchYear,
		&movie.Genre,
		&movie.Duration,
		&movie.Status,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func NewPGMovieRepository(pool *pgxpool.Pool) *PGMovieRepository {
	return &PGMovieRepository{
		pool,
	}
}
