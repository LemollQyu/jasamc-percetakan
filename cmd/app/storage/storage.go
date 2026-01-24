package storage

import (
	"context"
	"mime/multipart"
)

type Storage interface {
	UploadCategoryIcon(ctx context.Context, file *multipart.FileHeader) (string, error)
	UploadServiceIcon(ctx context.Context, file *multipart.FileHeader) (string, error)
	UploadServiceThumbnail(ctx context.Context, file *multipart.FileHeader) (string, error)
	UploadServiceGallery(ctx context.Context, file *multipart.FileHeader) (string, error)
	DeleteMediaByURL(ctx context.Context, url string) error
}
