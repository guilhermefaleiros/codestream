package service

import (
	"context"
	"fmt"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/application/gateway"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/application/model"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/entity"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/event"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/repository"
	"github.com/guilhermefaleiros/codestream/content-management-system/pkg"
)

type VideoService struct {
	storageGateway  gateway.StorageGateway
	videoRepository repository.VideoRepository
	eventMediator   *pkg.EventMediator
	baseFolder      string
}

func (v *VideoService) CreateVideo(ctx context.Context, input model.CreateVideoInput) (*model.CreateVideoOutput, error) {
	newVideo := entity.NewVideo(input.Title, input.Description, input.ContentType, input.OriginalFileName)
	storagePath := fmt.Sprintf("%s/%s%s", v.baseFolder, newVideo.ID, input.FileExtension)
	err := v.storageGateway.Upload(ctx, input.File, storagePath, input.ContentType)
	if err != nil {
		return nil, fmt.Errorf("unable to upload video: %w", err)
	}
	newVideo.Uploaded(storagePath, v.storageGateway.GetStorageLocation(), v.storageGateway.GetStorageProvider())
	err = v.videoRepository.Save(ctx, newVideo)
	if err != nil {
		return nil, fmt.Errorf("unable to save video: %w", err)
	}
	err = v.eventMediator.Dispatch(ctx, event.NewVideoUploadedEvent(newVideo))
	if err != nil {
		return nil, fmt.Errorf("unable to dispatch event: %w", err)
	}
	return &model.CreateVideoOutput{
		ID: newVideo.ID,
	}, nil
}

func NewVideoService(
	storageGateway gateway.StorageGateway,
	videoRepository repository.VideoRepository,
	eventMediator *pkg.EventMediator,
	baseFolder string,
) *VideoService {
	return &VideoService{
		storageGateway:  storageGateway,
		videoRepository: videoRepository,
		eventMediator:   eventMediator,
		baseFolder:      baseFolder,
	}
}
