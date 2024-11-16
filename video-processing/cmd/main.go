package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/guilhermefaleiros/codestream/video-processing/internal/application/service"
	"github.com/guilhermefaleiros/codestream/video-processing/internal/infrastructure/aws"
	appConfig "github.com/guilhermefaleiros/codestream/video-processing/internal/infrastructure/config"
)

func main() {
	ctx := context.Background()
	cfg := appConfig.LoadConfig()

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.Aws.AccessKeyID, cfg.Aws.SecretAccessKey, ""),
		),
		config.WithRegion(cfg.Aws.Region),
	)

	if err != nil {
		return
	}

	s3StorageGateway := aws.NewS3StorageGateway(awsCfg, cfg.Aws.S3.Bucket)
	transcodingConfig := service.TranscodingConfig{
		InputFolder:            "tmp",
		OutputBaseFolder:       "tmp",
		RemoteInputFolder:      cfg.Aws.S3.SourceFolder,
		RemoteOutputBaseFolder: cfg.Aws.S3.DestinationFolder,
	}
	transcodingService := service.NewTranscodingService(s3StorageGateway, transcodingConfig)

	videoID := "198d880f-84bb-4548-ac3a-2dd4d966544b"
	err = transcodingService.Execute(ctx, videoID)
}
