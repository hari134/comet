package pipelines

import (

	"github.com/hari134/comet/builder/implementations/stages"
	"github.com/hari134/comet/builder/pipeline"
)

var ReactViteNode20 = pipeline.NewPipeline().
		AddStage(stages.CopyProjectFilesToContainer()).
		AddStage(stages.ExtractProject()).
		AddStage(stages.InstallNpmDependencies()).
		AddStage(stages.NpmBuild()).
		AddStage(stages.CopyBuildFilesFromContainer())
