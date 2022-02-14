package tool

import (
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/wonderivan/logger"
)

// S3 - Amazon S3 storage
type S3 struct {
	Bucket            string
	RemotePath        string
	Region            string
	Access_key_id     string
	Secret_access_key string
	Endpoint          string
	Client            *s3manager.Uploader
}

// open connect
func (ctx *S3) Open() (err error) {

	sess, _ := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(ctx.Access_key_id, ctx.Secret_access_key, ""),
		Region:      aws.String(ctx.Region)},
	)
	ctx.Client = s3manager.NewUploader(sess)
	return
}

// close
func (ctx *S3) Close() {}

// upload
func (ctx *S3) Upload(srcFile, fileKey string) (err error) {
	f, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	remotePath := path.Join(ctx.RemotePath, fileKey)
	input := &s3manager.UploadInput{
		Bucket: aws.String(ctx.Bucket),
		Key:    aws.String(remotePath),
		Body:   f,
	}
	result, err := ctx.Client.Upload(input)
	if err != nil {
		return err
	}
	logger.Info(result.Location)
	return nil
}

// delete
func (ctx *S3) Delete(remotePath string) (err error) {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(ctx.Bucket),
		Key:    aws.String(remotePath),
	}
	_, err = ctx.Client.S3.DeleteObject(input)
	if err != nil {
		return err
	}
	return
}
