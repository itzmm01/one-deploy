package tool

import (
	"github.com/minio/minio-go/v6"
)

// Minio
type Minio struct {
	Bucket          string
	RemotePath      string
	Region          string
	Accesskeyid     string
	Secretaccesskey string
	Endpoint        string
	UseSSL          bool
}

// upload
func (ctx Minio) Upload(srcFile, fileKey string) (err error) {
	useSSL := false
	minioClient, err := minio.New(ctx.Endpoint, ctx.Accesskeyid, ctx.Secretaccesskey, useSSL)
	// 初使化minio client对象。
	if err != nil {
		return err
	}

	err = minioClient.MakeBucket(ctx.Bucket, ctx.Region)
	if err != nil {
		// 检查存储桶是否已经存在。
		exists, err := minioClient.BucketExists(ctx.Bucket)
		if err == nil && exists {
		} else {
			return err
		}
	}

	// 上传一个zip文件。
	objectName := fileKey
	filePath := srcFile
	contentType := "application/gzip"

	// 使用FPutObject上传一个zip文件。
	if _, err1 := minioClient.FPutObject(
		ctx.Bucket, objectName, filePath, minio.PutObjectOptions{ContentType: contentType},
	); err != nil {
		return err1
	}
	return
}

// delete
func (ctx Minio) Delete(srcFile string) (err error) {
	useSSL := false
	minioClient, err := minio.New(ctx.Endpoint, ctx.Accesskeyid, ctx.Secretaccesskey, useSSL)
	// 初使化minio client对象。
	if err != nil {
		return err
	}

	_, err = minioClient.StatObject(ctx.Bucket, srcFile, minio.StatObjectOptions{})
	if err != nil {
		return err
	}
	return minioClient.RemoveObject(ctx.Bucket, srcFile)
}
