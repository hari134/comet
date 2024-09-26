package pipelines

import (
	"github.com/hari134/comet/builder/pipeline"
)

var ReactViteNode20 pipeline.Pipeline

func InitializePipelines() {
	ReactViteNode20 = pipeline.NewSerialPipeline().
		AddStage(pipeline.NewFunctionStage(copyTarToContainer)).
		AddStage(pipeline.NewCommandStage("tar -xvf /app/full.tar -C /app")).
		AddStage(pipeline.NewCommandStage("cd /app && npm install")).
		AddStage(pipeline.NewCommandStage("cd /app && npm run build")).
		AddStage(pipeline.NewFunctionStage(copyDistFromContainer))
		// Add another stage to upload dist to s3, cdn
}


