package entity

import (
	"github.com/google/uuid"
	"time"
)

type Video struct {
	ID               string
	Title            string
	Description      string
	FilePath         string
	StorageLocation  string
	ContentType      string
	StorageProvider  string
	OriginalFileName string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Status           string
}

func (v *Video) Uploaded(filePath, storageLocation, storageProvider string) {
	v.FilePath = filePath
	v.StorageLocation = storageLocation
	v.StorageProvider = storageProvider
	v.Status = "uploaded"
}

func (v *Video) Processing() {
	v.Status = "processing"
}

func NewVideo(title, description, contentType, originalFileName string) *Video {
	return &Video{
		ID:               uuid.New().String(),
		Title:            title,
		ContentType:      contentType,
		Description:      description,
		OriginalFileName: originalFileName,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Status:           "created",
	}
}
