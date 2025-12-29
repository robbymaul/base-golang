package repositories

import (
	"context"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
)

func (rc *RepositoryContext) UploadFileRepository(ctx context.Context, assetPath string, file multipart.File, size int64) (info minio.UploadInfo, err error) {
	return rc.Minio.PutObject(ctx, rc.cfg.S3BucketName, assetPath, file, size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
}
