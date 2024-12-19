package cloud

import (
	"context"
	"mime/multipart"
	"time"
)

type CloudStorage interface {
	UploadContentFromMulipart(ctx context.Context, objectKey string, file multipart.File, contentType string) error
	UploadContentFromBytes(ctx context.Context, objectKey string, content []byte, contentType string) error
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	GetDownloadablePresignedURL(ctx context.Context, key string, duration time.Duration) (string, error)
	GetPresignedURL(ctx context.Context, key string, duration time.Duration) (string, error)
	GetContentByKey(ctx context.Context, objectKey string) ([]byte, error)
	GetMultipartFileByKey(ctx context.Context, objectKey string) (multipart.File, error)
	DeleteByKeys(ctx context.Context, objectKeys []string) error
}
