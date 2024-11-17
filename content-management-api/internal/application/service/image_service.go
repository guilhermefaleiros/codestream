package service

import (
	"context"
	"fmt"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/application/gateway"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/application/model"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/entity"
	"github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/repository"
	"github.com/guilhermefaleiros/codestream/content-management-system/pkg"
)

type ImageService struct {
	imageRepository repository.ImageRepository
	storageGateway  gateway.StorageGateway
}

func (i *ImageService) CreateImage(ctx context.Context, input model.CreateImageInput) (*model.CreateImageOutput, error) {
	existentImageType, _ := i.imageRepository.FindByMovieIDAndType(ctx, input.MovieID, input.Type)

	if existentImageType != nil {
		return nil, fmt.Errorf("image already exists")
	}

	contentType := pkg.GetFileContentType(input.File)
	fileExtension, _ := pkg.GetMimeToExtension(contentType)

	newImage := entity.NewImage(contentType, input.OriginalFileName, input.MovieID, input.Type)
	filePath := fmt.Sprintf("%s/%s/%s%s", input.MovieID, "images", newImage.ID, fileExtension)

	err := i.storageGateway.Upload(ctx, input.File, filePath, contentType)
	if err != nil {
		return nil, fmt.Errorf("unable to upload image: %w", err)
	}

	newImage.Uploaded(filePath, i.storageGateway.GetStorageLocation(), i.storageGateway.GetStorageProvider())
	err = i.imageRepository.Save(ctx, newImage)
	if err != nil {
		return nil, fmt.Errorf("unable to save image: %w", err)
	}
	return &model.CreateImageOutput{
		ID: newImage.ID,
	}, nil
}

func NewImageService(
	imageRepository repository.ImageRepository,
	storageGateway gateway.StorageGateway,
) *ImageService {
	return &ImageService{
		imageRepository: imageRepository,
		storageGateway:  storageGateway,
	}
}
