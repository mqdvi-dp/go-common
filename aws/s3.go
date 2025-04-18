package aws

import (
	"bytes"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/logger"
)

type awsS3 struct {
	bucketName string
	clientS3   *s3.Client
}

type AwsS3 interface {
	// PutObject put object to s3
	PutObject(ctx context.Context, objectName string, file []byte, contentTypes ...string) (*s3.PutObjectOutput, error)
	// GetSignURL get presign url for object
	GetSignURL(ctx context.Context, objectName string, durations ...time.Duration) (string, error)
	// DeleteObject delete object from s3
	DeleteObject(ctx context.Context, objectName string) (*s3.DeleteObjectOutput, error)
	// GetObject get object from s3
	GetObject(ctx context.Context, objectName string) (*s3.GetObjectOutput, error)
}

// newS3 is a function that returns a new instance of awsS3
func newS3(cfg aws.Config, bucketName string) AwsS3 {
	client := s3.NewFromConfig(cfg)
	if client == nil {
		logger.Log.Fatalf("failed to create s3 client")
	}

	return &awsS3{
		clientS3:   client,
		bucketName: bucketName,
	}
}

func (as *awsS3) GetSignURL(ctx context.Context, objectName string, durations ...time.Duration) (url string, err error) {
	duration := env.GetDuration("DEFAULT_DURATION_S3_PRESIGN_URL", 1*time.Hour)
	if len(durations) > 0 {
		duration = durations[0]
	}

	signClient := s3.NewPresignClient(as.clientS3, s3.WithPresignExpires(duration))
	result, err := signClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(as.bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		logger.Log.Errorf(ctx, "failed to get object from s3: %v", err)
		return
	}

	url = result.URL
	return
}

func (as *awsS3) PutObject(ctx context.Context, objectName string, file []byte, contentTypes ...string) (*s3.PutObjectOutput, error) {
	contentType := "application/octet-stream"
	if len(contentTypes) > 0 {
		contentType = contentTypes[0]
	}

	res, err := as.clientS3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(as.bucketName),
		Key:         aws.String(objectName),
		ContentType: aws.String(contentType),
		Body:        bytes.NewReader(file),
	})
	if err != nil {
		logger.Log.Errorf(ctx, "failed to put object to s3: %v", err)
		return nil, err
	}

	return res, nil
}

func (as *awsS3) DeleteObject(ctx context.Context, objectName string) (*s3.DeleteObjectOutput, error) {
	res, err := as.clientS3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(as.bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		logger.Log.Errorf(ctx, "failed to delete object from s3: %v", err)
		return nil, err
	}

	return res, nil
}

func (as *awsS3) GetObject(ctx context.Context, objectName string) (*s3.GetObjectOutput, error) {
	res, err := as.clientS3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(as.bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		logger.Log.Errorf(ctx, "failed to get object from s3: %v", err)
		return nil, err
	}

	return res, nil
}
