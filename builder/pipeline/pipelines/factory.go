package pipelines

import (
	"errors"
	"github.com/hari134/comet/builder/pipeline"
)

func PipelineFactory(buildType string) (pipeline.Pipeline, error){
	switch buildType{
	case "ReactViteNode20":
		return ReactViteNode20,nil
	default:
		return nil,errors.New("no such pipeline exists")
	}
}

