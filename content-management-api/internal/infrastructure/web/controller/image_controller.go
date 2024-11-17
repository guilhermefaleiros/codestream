package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/application/model"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/application/service"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/entity"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/infrastructure/web"
	"net/http"
)

type ImageController struct {
	r            chi.Router
	imageService *service.ImageService
	maxFileSize  int64
}

func (ic *ImageController) SetupRoutes() {
	ic.r.Post("/v1/movie/{movieID}/image", ic.UploadImage)
}

func (ic *ImageController) UploadImage(w http.ResponseWriter, r *http.Request) {
	movieID := chi.URLParam(r, "movieID")

	err := r.ParseMultipartForm(ic.maxFileSize)
	if err != nil {
		web.JSONError(w, "file too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		web.JSONError(w, "missing file", http.StatusBadRequest)
		return
	}

	fileBytes, err := web.FileToBytes(file)
	if err != nil {
		web.JSONError(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	input := model.CreateImageInput{
		MovieID:          movieID,
		Type:             entity.ImageType(r.FormValue("type")),
		File:             fileBytes,
		OriginalFileName: header.Filename,
	}

	image, err := ic.imageService.CreateImage(r.Context(), input)
	if err != nil {
		web.JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	web.JSON(w, image, http.StatusCreated)
}

func NewImageController(r chi.Router, imageService *service.ImageService, maxFileSize int64) *ImageController {
	return &ImageController{
		r:            r,
		imageService: imageService,
		maxFileSize:  maxFileSize,
	}
}
