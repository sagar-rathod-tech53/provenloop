package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	aws_credentials "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	app_config "github.com/sagar-rathod-devops/do-host-network-backend/config"
)

type S3Uploader struct {
	Client     *s3.Client
	BucketName string
	Region     string
}

func NewS3Uploader(cfg app_config.Config) (*S3Uploader, error) {
	awsCfg, err := aws_config.LoadDefaultConfig(context.TODO(),
		aws_config.WithRegion(cfg.AWS_REGION),
		aws_config.WithCredentialsProvider(aws_credentials.NewStaticCredentialsProvider(
			cfg.AWS_ACCESS_KEY_ID,
			cfg.AWS_SECRET_ACCESS_KEY,
			"",
		)),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)

	return &S3Uploader{
		Client:     client,
		BucketName: cfg.AWS_BUCKET_NAME,
		Region:     cfg.AWS_REGION,
	}, nil
}

func (u *S3Uploader) UploadFile(file multipart.File, fileHeader *multipart.FileHeader, key string) (string, error) {
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

	_, err := u.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(u.BucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.BucketName, u.Region, key)
	return url, nil
}
