package event

import (
	"github.com/google/uuid"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/entity"
	"time"
)

type VideoUploadedEventPayload struct {
	ID              string
	FilePath        string
	StorageLocation string
	StorageProvider string
	CreatedAt       time.Time
}

type VideoUploadedEvent struct {
	ID        string
	Payload   VideoUploadedEventPayload
	CreatedAt time.Time
}

func (e VideoUploadedEvent) GetID() string {
	return e.ID
}

func (e VideoUploadedEvent) GetCreatedAt() time.Time {
	return e.CreatedAt
}

func (e VideoUploadedEvent) GetPayload() interface{} {
	return e.Payload
}

func (e VideoUploadedEvent) GetType() string {
	return "VideoUploaded"
}

func NewVideoUploadedEvent(video *entity.Video) *VideoUploadedEvent {
	return &VideoUploadedEvent{
		ID: uuid.New().String(),
		Payload: VideoUploadedEventPayload{
			ID:              video.ID,
			FilePath:        video.FilePath,
			StorageLocation: video.StorageLocation,
			StorageProvider: video.StorageProvider,
			CreatedAt:       video.CreatedAt,
		},
		CreatedAt: time.Now(),
	}
}
