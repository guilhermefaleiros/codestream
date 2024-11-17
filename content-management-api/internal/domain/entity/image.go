package entity

import (
	"github.com/google/uuid"
	"time"
)

type ImageType string

const (
	ImageThumbnail ImageType = "thumbnail"
	ImageCover     ImageType = "cover"
	ImagePoster    ImageType = "poster"
)

type Image struct {
	ID               string
	FilePath         string
	StorageLocation  string
	ContentType      string
	StorageProvider  string
	OriginalFileName string
	MovieID          string
	Type             ImageType
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (i *Image) Uploaded(filePath, storageLocation, storageProvider string) {
	i.FilePath = filePath
	i.StorageLocation = storageLocation
	i.StorageProvider = storageProvider
}

func NewImage(contentType, originalFileName, movieID string, imageType ImageType) *Image {
	return &Image{
		ID:               uuid.New().String(),
		ContentType:      contentType,
		OriginalFileName: originalFileName,
		MovieID:          movieID,
		Type:             imageType,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}
