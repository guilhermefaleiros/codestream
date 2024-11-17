package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/application/model"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/application/service"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/infrastructure/web"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/infrastructure/web/dto"
	"net/http"
)

type MovieController struct {
	r            chi.Router
	movieService *service.MovieService
}

func (mc *MovieController) SetupRoutes() {
	mc.r.Post("/v1/movies", mc.CreateEmptyMovie)
}

func (mc *MovieController) CreateEmptyMovie(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateEmptyMovieRequestDTO

	err := web.DeserializeRequestBody(r, &request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	input := model.CreateEmptyMovieInput{
		Title:       request.Title,
		Description: request.Description,
		Genre:       request.Genre,
		LaunchYear:  request.LaunchYear,
		Duration:    request.Duration,
	}
	output, err := mc.movieService.CreateEmptyMovie(r.Context(), input)

	if err != nil {
		web.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	web.JSON(w, output, http.StatusCreated)
}

func NewMovieController(r chi.Router, movieService *service.MovieService) *MovieController {
	return &MovieController{
		r:            r,
		movieService: movieService,
	}
}
