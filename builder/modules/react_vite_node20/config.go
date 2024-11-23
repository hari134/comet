package reactvitenode20

import (
	"bytes"
	"context"

	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/builder/relay"
	"github.com/hari134/comet/core/cdn"
	"github.com/hari134/comet/core/dns"
	"github.com/hari134/comet/core/storage"
)

// PipelineOptions holds the configuration for the entire pipeline.
type PipelineOptions struct {
	BuildContainer       container.BuildContainer
	ProjectStorageConfig ProjectStorageConfig
	ProjectFileData      ProjectFileData
	BuildOutputFilesData BuildOutputFilesData
	Store                storage.Store
	CDNProvider          cdn.CDNProvider
	DNSProvider          dns.DNSProvider
	StreamConfig         StreamConfig
	DeploymentConfig     DeploymentConfig
	Context              context.Context
}

// ProjectStorageConfig handles storage-related configurations for the project.
type ProjectStorageConfig struct {
	ProjectFilesBucket string
	BuildFilesBucket   string
	ProjectFilesKey    string
	BuildFilesKey      string
}

// ProjectFileData encapsulates details about the project's input files.
type ProjectFileData struct {
	DirName        string
	ProjectTarFile *bytes.Buffer
}

// BuildOutputFilesData encapsulates details about the project's build output files.
type BuildOutputFilesData struct {
	DirName        string
	ProjectTarFile *bytes.Buffer
}

// StreamConfig handles streaming-related options and data.
type StreamConfig struct {
	Enabled bool
	Output  chan relay.StreamData
}

// DeploymentConfig handles deployment-related configurations.
type DeploymentConfig struct {
	Subdomain        string
	MainDomain       string
	CloudFrontConfig CloudFrontConfig
	Route53Config    Route53Config
}

// CloudFrontConfig encapsulates configurations for CloudFront.
type CloudFrontConfig struct {
	CertificateARN string
	OAI            string // Origin Access Identity
	Region         string
}

// Route53Config encapsulates configurations for DNS management.
type Route53Config struct {
	HostedZoneID string
}

// Enable enables streaming in the StreamConfig.
func (s *StreamConfig) Enable() {
	s.Enabled = true
}

// IsEnabled checks if streaming is enabled in the StreamConfig.
func (s *StreamConfig) IsEnabled() bool {
	return s.Enabled
}

func (s *PipelineOptions) IsStreamingEnabled() bool {
	return s.StreamConfig.IsEnabled()
}

// NewPipelineOptions initializes and returns a PipelineOptions instance.
func NewPipelineOptions(
	buildContainer container.BuildContainer,
	store storage.Store,
	cdnProvider cdn.CDNProvider,
	dnsProvider dns.DNSProvider,
	projectStorageConfig ProjectStorageConfig,
	deploymentConfig DeploymentConfig,
) *PipelineOptions {
	return &PipelineOptions{
		BuildContainer:       buildContainer,
		Store:                store,
		CDNProvider:          cdnProvider,
		DNSProvider:          dnsProvider,
		ProjectStorageConfig: projectStorageConfig,
		DeploymentConfig:     deploymentConfig,
		StreamConfig: StreamConfig{
			Enabled: false,
			Output:  make(chan relay.StreamData),
		},
	}
}
