package fai

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type s3options struct {
	client          *minio.Client
	bucketName      string
	address         string
	accessKeyID     string
	secretAccessKey string
	token           string
	useSSL          bool
}

type s3Option func(*s3options)

func WithS3Address(address string) s3Option {
	return func(opts *s3options) {
		opts.address = address
	}
}

func WithS3AccessKeyID(accessKeyID string) s3Option {
	return func(opts *s3options) {
		opts.accessKeyID = accessKeyID
	}
}

func WithS3SecretAccessKey(secretAccessKey string) s3Option {
	return func(opts *s3options) {
		opts.secretAccessKey = secretAccessKey
	}
}

func WithS3Token(token string) s3Option {
	return func(opts *s3options) {
		opts.token = token
	}
}

func WithS3BucketName(bucket string) s3Option {
	return func(opts *s3options) {
		opts.bucketName = bucket
	}
}

func NewS3Uploader(options ...s3Option) (*s3options, error) {
	o := new(s3options)

	for _, opt := range options {
		opt(o)
	}
	if o.bucketName == "" {
		return nil, errors.New("s3 bucket name is required")
	}

	s3credentials := credentials.NewStaticV4(o.accessKeyID, o.secretAccessKey, o.token)
	o.useSSL = true
	s3Client, err := minio.New(o.address, &minio.Options{
		Creds:  s3credentials,
		Secure: o.useSSL,
	})
	if err != nil {
		return nil, err
	}
	o.client = s3Client

	return o, nil
}

// Upload uploads the file at filePath to S3.
func (f *s3options) Upload(filePath string) (minio.UploadInfo, error) {
	ctx := context.Background()
	bucketName := f.bucketName
	objectName := filepath.Base(filePath)
	options := minio.PutObjectOptions{}

	return f.client.FPutObject(ctx, bucketName, objectName, filePath, options)
}