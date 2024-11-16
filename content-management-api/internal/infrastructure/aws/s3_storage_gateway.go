package aws

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

type S3StorageGateway struct {
	client     *s3.Client
	bucketName string
}

func (s *S3StorageGateway) Upload(ctx context.Context, file []byte, path, contentType string) error {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(path),
		Body:        bytes.NewReader(file),
		ContentType: aws.String(contentType),
	}
	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	log.Printf("File uploaded successfully to bucket %s with key %s", s.bucketName, path)
	return nil
}

func (s *S3StorageGateway) GetStorageLocation() string {
	return s.bucketName
}

func (s *S3StorageGateway) GetStorageProvider() string {
	return "AWS_S3"
}

func NewS3StorageGateway(config aws.Config, bucketName string) *S3StorageGateway {
	client := s3.NewFromConfig(config)
	return &S3StorageGateway{
		client, bucketName,
	}
}
