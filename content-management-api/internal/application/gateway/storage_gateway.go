package gateway

import "context"

type StorageGateway interface {
	Upload(ctx context.Context, file []byte, path, contentType string) error
	GetStorageLocation() string
	GetStorageProvider() string
}
