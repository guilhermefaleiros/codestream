package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/application/model"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/application/service"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/entity"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/infrastructure/web"
	"net/http"
)

type VideoController struct {
	r            chi.Router
	videoService *service.VideoService
	maxFileSize  int64
}

func (vc *VideoController) SetupRoutes() {
	vc.r.Post("/v1/movie/{movieID}/video", vc.UploadVideo)
}

func (vc *VideoController) UploadVideo(w http.ResponseWriter, r *http.Request) {
	movieID := chi.URLParam(r, "movieID")

	err := r.ParseMultipartForm(vc.maxFileSize)
	if err != nil {
		web.JSONError(w, "file too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		web.JSONError(w, "missing file", http.StatusBadRequest)
		return
	}

	fileBytes, err := web.FileToBytes(file)
	if err != nil {
		web.JSON(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	input := model.CreateVideoInput{
		MovieID:          movieID,
		Type:             entity.VideoType(r.FormValue("type")),
		File:             fileBytes,
		OriginalFileName: header.Filename,
	}

	output, err := vc.videoService.CreateVideo(r.Context(), input)

	if err != nil {
		web.JSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	web.JSON(w, output, http.StatusCreated)
}

func NewVideoController(r chi.Router, videoService *service.VideoService, maxFileSize int64) *VideoController {
	return &VideoController{
		r:            r,
		videoService: videoService,
		maxFileSize:  maxFileSize,
	}
}
