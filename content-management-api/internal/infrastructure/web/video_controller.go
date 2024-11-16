package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/application/model"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/application/service"
	"net/http"
)

type VideoController struct {
	r            chi.Router
	videoService *service.VideoService
	maxFileSize  int64
}

func (vc *VideoController) SetupRoutes() {
	vc.r.Post("/v1/videos", vc.Upload)
}

func (vc *VideoController) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(vc.maxFileSize)
	if err != nil {
		JSONError(w, "file too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		JSONError(w, "missing file", http.StatusBadRequest)
		return
	}

	fileBytes, err := FileToBytes(file)
	if err != nil {
		JSON(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	contentType := GetFileContentType(fileBytes)
	extension, err := GetMimeToExtension(contentType)

	if err != nil {
		JSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	input := model.CreateVideoInput{
		Title:            r.FormValue("title"),
		Description:      r.FormValue("description"),
		File:             fileBytes,
		ContentType:      contentType,
		FileExtension:    extension,
		OriginalFileName: header.Filename,
	}

	output, err := vc.videoService.CreateVideo(r.Context(), input)

	if err != nil {
		JSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	JSON(w, output, http.StatusCreated)
}

func NewVideoController(r chi.Router, videoService *service.VideoService, maxFileSize int64) *VideoController {
	return &VideoController{
		r:            r,
		videoService: videoService,
		maxFileSize:  maxFileSize,
	}
}
