package dns

import (
	"context"
)

// DNSProvider defines a generic interface for DNS operations
type DNSProvider interface {
	CreateOrUpdateRecord(ctx context.Context, recordName, target, recordType string) error
	DeleteRecord(ctx context.Context, recordName, recordType string) error
}
