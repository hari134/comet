package reactvitenode20

import "github.com/hari134/comet/builder/pipeline"

var ReactViteNode20Pipeline = pipeline.NewPipeline().
	AddStage(PullProjectFiles()).
	AddStage(CopyProjectFilesToContainer()).
	// AddStage(ExtractProject()).
	AddStage(InstallNpmDependencies()).
	AddStage(NpmBuild()).
	AddStage(CopyBuildFilesFromContainer()).
	AddStage(UploadBuildFilesToS3())

	// Stage to upload to s3 and cdn
	// Stage to setup routing in Route 53
	// Stage to add data to db
