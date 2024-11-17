package service

import (
	"context"
	"fmt"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/application/model"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/entity"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/repository"
)

type MovieService struct {
	movieRepository repository.MovieRepository
}

func (m *MovieService) CreateEmptyMovie(ctx context.Context, input model.CreateEmptyMovieInput) (*model.CreateMovieOutput, error) {
	newMovie := entity.NewMovie(input.Title, input.Description, input.Genre, input.LaunchYear, input.Duration)
	err := m.movieRepository.Save(ctx, newMovie)
	if err != nil {
		return nil, fmt.Errorf("unable to save movie: %w", err)
	}
	return &model.CreateMovieOutput{
		ID: newMovie.ID,
	}, nil
}

func NewMovieService(
	movieRepository repository.MovieRepository,
) *MovieService {
	return &MovieService{
		movieRepository: movieRepository,
	}
}
