package entity

import (
	"github.com/google/uuid"
	"time"
)

type VideoType string

const (
	VideoTypeTrailer VideoType = "trailer"
	VideoTypeMovie   VideoType = "movie"
)

type VideoStatus string

const (
	VideoStatusCreated    VideoStatus = "created"
	VideoStatusUploaded   VideoStatus = "uploaded"
	VideoStatusProcessing VideoStatus = "processing"
)

type Video struct {
	ID               string
	FilePath         string
	StorageLocation  string
	ContentType      string
	StorageProvider  string
	OriginalFileName string
	MovieID          string
	Type             VideoType
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Status           VideoStatus
}

func (v *Video) Uploaded(filePath, storageLocation, storageProvider string) {
	v.FilePath = filePath
	v.StorageLocation = storageLocation
	v.StorageProvider = storageProvider
	v.Status = VideoStatusUploaded
}

func (v *Video) Processing() {
	v.Status = VideoStatusProcessing
}

func NewVideo(contentType, originalFileName, movieID string, videoType VideoType) *Video {
	return &Video{
		ID:               uuid.New().String(),
		ContentType:      contentType,
		OriginalFileName: originalFileName,
		MovieID:          movieID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Type:             videoType,
		Status:           VideoStatusCreated,
	}
}
