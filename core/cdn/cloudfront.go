package cdn

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

type CloudfrontConfig struct {
	AccessKey string
	SecretKey string
	Region    string
}

// CloudFrontCDN implements the CDN interface for AWS CloudFront
type CloudFrontCDN struct {
	client *cloudfront.Client
}

// DistributionConfig holds the configuration for creating a CloudFront distribution
type DistributionConfig struct {
	OriginDomain      string
	DomainName        string
	CertificateARN    string
	EnableLogging     bool
	LogBucket         string
	DefaultRootObject string
	OAI               string // Origin Access Identity
}

func NewCloudFrontCDN(config CloudfrontConfig) *CloudFrontCDN {
	cfg := aws.Config{
		Region: config.Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			config.AccessKey, config.SecretKey, "",
		),
	}

	client := cloudfront.New(cloudfront.Options{
		Credentials: cfg.Credentials,
		Region:      cfg.Region,
	})
	return &CloudFrontCDN{client: client}
}

// CreateDistribution creates a new CloudFront distribution with OAI
func (cdn *CloudFrontCDN) CreateDistribution(ctx context.Context, config DistributionConfig) (string, error) {
	input := &cloudfront.CreateDistributionInput{
		DistributionConfig: &types.DistributionConfig{
			CallerReference: aws.String(fmt.Sprintf("dist-%s", config.DomainName)),
			Origins: &types.Origins{
				Items: []types.Origin{
					{
						Id:         aws.String("S3Origin"),
						DomainName: aws.String(config.OriginDomain),
						S3OriginConfig: &types.S3OriginConfig{
							OriginAccessIdentity: aws.String(config.OAI), // Set the OAI here
						},
					},
				},
				Quantity: aws.Int32(1),
			},
			DefaultCacheBehavior: &types.DefaultCacheBehavior{
				TargetOriginId:       aws.String("S3Origin"),
				ViewerProtocolPolicy: types.ViewerProtocolPolicyRedirectToHttps,
				Compress:             aws.Bool(true),
			},
			Enabled: aws.Bool(true),
			ViewerCertificate: &types.ViewerCertificate{
				ACMCertificateArn: aws.String(config.CertificateARN),
				SSLSupportMethod:  types.SSLSupportMethodSniOnly,
			},
			DefaultRootObject: aws.String(config.DefaultRootObject),
			Logging: &types.LoggingConfig{
				Enabled: aws.Bool(config.EnableLogging),
				Bucket:  aws.String(config.LogBucket),
			},
		},
	}

	result, err := cdn.client.CreateDistribution(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to create CloudFront distribution: %w", err)
	}

	return *result.Distribution.DomainName, nil
}

// DeleteDistribution disables and deletes an existing CloudFront distribution
func (cdn *CloudFrontCDN) DeleteDistribution(ctx context.Context, distributionID string) error {
	// Disable the distribution before deleting it
	_, err := cdn.client.UpdateDistribution(ctx, &cloudfront.UpdateDistributionInput{
		Id: aws.String(distributionID),
		DistributionConfig: &types.DistributionConfig{
			Enabled: aws.Bool(false),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to disable CloudFront distribution: %w", err)
	}

	// Delete the distribution
	_, err = cdn.client.DeleteDistribution(ctx, &cloudfront.DeleteDistributionInput{
		Id: aws.String(distributionID),
	})
	if err != nil {
		return fmt.Errorf("failed to delete CloudFront distribution: %w", err)
	}

	return nil
}
