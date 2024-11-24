package cdn

import "context"

// CDNProvider defines an interface for content delivery network operations
type CDNProvider interface {
	CreateDistribution(ctx context.Context, config DistributionConfig) (string, error)
	DeleteDistribution(ctx context.Context, distributionID string) error
}
