package dns

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// Route53Config holds the configuration for Route 53
type Route53Config struct {
	AccessKey    string
	SecretKey    string
	Region       string
	HostedZoneID string
}

// Route53DNS implements DNSProvider for AWS Route 53
type Route53DNS struct {
	client       *route53.Client
	hostedZoneID string
}

// NewRoute53DNS initializes a new Route53DNS instance with explicit configuration
func NewRoute53DNS(config Route53Config) (*Route53DNS, error) {
	awsCfg := aws.Config{
		Region: config.Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			config.AccessKey,
			config.SecretKey,
			"",
		),
	}

	client := route53.New(route53.Options{
		Credentials: awsCfg.Credentials,
		Region:      awsCfg.Region,
	})

	return &Route53DNS{
		client:       client,
		hostedZoneID: config.HostedZoneID,
	}, nil
}

// CreateOrUpdateRecord creates or updates a DNS record in Route 53
func (r *Route53DNS) CreateOrUpdateRecord(ctx context.Context, recordName, target, recordType string) error {
	_, err := r.client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(r.hostedZoneID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(recordName),
						Type: types.RRType(recordType),
						TTL:  aws.Int64(300),
						ResourceRecords: []types.ResourceRecord{
							{Value: aws.String(target)},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create or update DNS record: %w", err)
	}
	return nil
}

// DeleteRecord deletes a DNS record in Route 53
func (r *Route53DNS) DeleteRecord(ctx context.Context, recordName, recordType string) error {
	_, err := r.client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(r.hostedZoneID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionDelete,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(recordName),
						Type: types.RRType(recordType),
						TTL:  aws.Int64(300),
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete DNS record: %w", err)
	}
	return nil
}
