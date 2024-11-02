package pipelines

import (

	"github.com/hari134/comet/builder/implementations/stages"
	"github.com/hari134/comet/builder/pipeline"
)

var ReactViteNode20 = pipeline.NewPipeline().
// Add pull from aws stage as well
		AddStage(stages.PullProjectFiles()).
		AddStage(stages.CopyProjectFilesToContainer()).
		AddStage(stages.ExtractProject()).
		AddStage(stages.InstallNpmDependencies()).
		AddStage(stages.NpmBuild()).
		AddStage(stages.CopyBuildFilesFromContainer())
		// Stage to upload to s3 and cdn
		// Stage to setup routing in Route 53
		// Stage to add data to db

