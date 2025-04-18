package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AWS interface {
	// S3 creates a new S3 client
	S3(bucketName string) AwsS3
	// TODO: Add other AWS services here
}

type awsConfig struct {
	config aws.Config
}

// New is a function that returns a new instance of awsConfig
func New() AWS {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}

	return &awsConfig{
		config: cfg,
	}
}

func (ac *awsConfig) S3(bucketName string) AwsS3 {
	return newS3(ac.config, bucketName)
}
