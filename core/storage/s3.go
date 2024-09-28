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
	BucketName      string
	Region          string
}

type S3Store struct {
	client *s3.S3
	bucket string
}

func (s3Store S3Store) Get(ctx context.Context,key string) ([]byte, error) {
	resp, err := s3Store.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3Store.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (s3Store S3Store) Put(ctx context.Context,fileData []byte, key string) error {
	_, err := s3Store.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Store.bucket),
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
	s3Store := &S3Store{client: client, bucket: awsConfig.BucketName}
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

