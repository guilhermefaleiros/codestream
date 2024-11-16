package gateway

import "context"

type StorageGateway interface {
	Upload(ctx context.Context, file []byte, path, contentType string) error
	Download(ctx context.Context, path string) ([]byte, error)
	GetStorageLocation() string
	GetStorageProvider() string
}
