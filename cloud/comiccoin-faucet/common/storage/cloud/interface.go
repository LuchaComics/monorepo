package cloud

import (
	"context"
	"mime/multipart"
	"time"
)

type CloudStorage interface {
	UploadContent(ctx context.Context, objectKey string, content []byte) error
	UploadContentFromMulipart(ctx context.Context, objectKey string, file multipart.File) error
	UploadContentFromBytes(ctx context.Context, objectKey string, content []byte) error
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	GetDownloadablePresignedURL(ctx context.Context, key string, duration time.Duration) (string, error)
	GetPresignedURL(ctx context.Context, key string, duration time.Duration) (string, error)
	GetContentByKey(ctx context.Context, objectKey string) ([]byte, error)
	GetMultipartFileByKey(ctx context.Context, objectKey string) (multipart.File, error)
	DeleteByKeys(ctx context.Context, key []string) error
}
