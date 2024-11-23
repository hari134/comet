package builder

import (
	"context"
	"fmt"
	"github.com/hari134/comet/core/cdn"
	"github.com/hari134/comet/core/dns"

	"github.com/hari134/comet/builder/container"
	"github.com/hari134/comet/builder/modules"
	reactvitenode20 "github.com/hari134/comet/builder/modules/react_vite_node20"
	"github.com/hari134/comet/builder/relay"
	"github.com/hari134/comet/core/storage"
)

type Builder struct {
	Store            storage.Store
	CDNProvider      cdn.CDNProvider
	DNSProvider      dns.DNSProvider
	ContainerManager container.ContainerManager
	PipelineFactory  modules.PipelineFactory
}

func NewBuilder(store storage.Store, containerManager container.ContainerManager, pipelineFactory modules.PipelineFactory) *Builder {
	return &Builder{
		Store:            store,
		ContainerManager: containerManager,
		PipelineFactory:  pipelineFactory,
	}
}

type ProjectDeploymentConfig struct {
	ProjectStorageKey    string
	ProjectStorageBucket string
	BuildFilesBucket     string
	BuildEnvType         string
	OriginDomain         string
	SubDomain            string
	CloudFrontOAI        string
	CloudFrontRegion     string
	CertificateARN       string
	Route53HostedZoneID  string
}

func (builder *Builder) DeployProject(config ProjectDeploymentConfig) error {
	buildEnvType := config.BuildEnvType
	buildPipeline, err := builder.PipelineFactory.Get(buildEnvType)

	if err != nil {
		return err
	}
	switch buildEnvType {
	case "reactvitenode20":
		buildContainer, err := builder.ContainerManager.NewBuildContainer(buildEnvType)
		if err != nil {
			return err
		}

		cfg := reactvitenode20.PipelineOptions{
			BuildContainer: buildContainer,

			ProjectStorageConfig: reactvitenode20.ProjectStorageConfig{
				ProjectFilesKey:    config.ProjectStorageKey,
				ProjectFilesBucket: config.ProjectStorageBucket,
				BuildFilesBucket:   config.BuildFilesBucket,
			},

			Store:       builder.Store,
			CDNProvider: builder.CDNProvider,
			DNSProvider: builder.DNSProvider,

			StreamConfig: reactvitenode20.StreamConfig{
				Enabled: true,
				Output:  make(chan relay.StreamData),
			},

			DeploymentConfig: reactvitenode20.DeploymentConfig{
				Subdomain:  config.SubDomain,
				MainDomain: config.OriginDomain,
				CloudFrontConfig: reactvitenode20.CloudFrontConfig{
					CertificateARN: config.CertificateARN,
					OAI:            config.CloudFrontOAI,
					Region:         config.CloudFrontRegion,
				},
				Route53Config: reactvitenode20.Route53Config{
					HostedZoneID: config.Route53HostedZoneID,
				},
			},

			Context: context.Background(), // Use custom context with deadline ,cancelling
		}
		if err := buildPipeline.Run(&cfg); err != nil {
			return err
		}
	default:
		return fmt.Errorf("no such build environment type %v", buildEnvType)
	}
	return nil
}
