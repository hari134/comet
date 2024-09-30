package storage

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AWSCredentials struct {
	AccessKey       string
	SecretAccessKey string
	Region          string
}

type S3Store struct {
	client *s3.S3
}

func (s3Store S3Store) Get(ctx context.Context, bucket string, key string) (*bytes.Buffer, error) {
	resp, err := s3Store.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	// Ensure the response body is closed after reading
	defer resp.Body.Close()

	// Read the object body into a byte slice
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(data)

	return buffer, nil
}

func (s3Store S3Store) Put(ctx context.Context, fileData []byte, bucket, key string) error {
	_, err := s3Store.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fileData),
	})
	if err != nil {
		return err
	}
	return nil
}

func NewS3Store(awsConfig AWSCredentials) (*S3Store, error) {
	sess, err := newAwsSession(awsConfig)
	if err != nil {
		return nil, err
	}
	client := s3.New(sess)
	s3Store := &S3Store{client: client}
	return s3Store, nil
}

func newAwsSession(config AWSCredentials) (*session.Session, error) {
	awsRegion := config.Region
	accessKey := config.AccessKey
	secretKey := config.SecretAccessKey

	awsConfig := &aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	}
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	return sess, nil
}
