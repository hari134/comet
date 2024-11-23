package storage

import (
	"bytes"
	"context"
	"io"

	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AWSCredentials struct {
	AccessKey       string
	SecretAccessKey string
	Region          string
}

type S3Store struct {
	client *s3.Client
}

// Get retrieves an object from the S3 bucket and returns it as a bytes.Buffer
func (s3Store *S3Store) Get(ctx context.Context, bucket string, key string) (*bytes.Buffer, error) {
	resp, err := s3Store.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		slog.Debug("Failed to retrieve object from S3", "bucket", bucket, "key", key, "error", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Debug("Failed to read object data", "bucket", bucket, "key", key, "error", err.Error())
		return nil, err
	}

	return bytes.NewBuffer(data), nil
}

// Put uploads a file to the specified S3 bucket
func (s3Store *S3Store) Put(ctx context.Context, fileData *bytes.Buffer, bucket, key string) error {
	_, err := s3Store.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fileData.Bytes()),
	})

	if err != nil {
		slog.Debug("Failed to upload to S3", "bucket", bucket, "key", key, "error", err.Error())
		return err
	}
	return nil
}

// NewS3Store initializes and returns a new S3Store instance
func NewS3Store(awsConfig AWSCredentials) (*S3Store, error) {
	// Load AWS configuration with static credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsConfig.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(awsConfig.AccessKey, awsConfig.SecretAccessKey, ""),
		),
	)
	if err != nil {
		slog.Debug("Failed to load AWS configuration", "error", err.Error())
		return nil, err
	}

	// Create the S3 client
	client := s3.NewFromConfig(cfg)
	return &S3Store{client: client}, nil
}
