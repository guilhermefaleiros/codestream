package model

import "github.com/guilhermefaleiros/codestream/content-management-system/internal/domain/entity"

type CreateVideoInput struct {
	MovieID          string
	Type             entity.VideoType
	OriginalFileName string
	File             []byte
}

type CreateImageInput struct {
	MovieID          string
	Type             entity.ImageType
	OriginalFileName string
	File             []byte
}

type CreateEmptyMovieInput struct {
	Title       string
	Description string
	Genre       string
	LaunchYear  int
	Duration    int
}
