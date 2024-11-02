package stages

import (
	"context"
	"log/slog"

	"github.com/hari134/comet/builder/dto"
	"github.com/hari134/comet/builder/pipeline"
	"github.com/hari134/comet/core/storage"
)

func PullProjectFiles() pipeline.Stage {
	return pipeline.Stage{
		Name: "Pull Project Files From Storage",
		Execute: func(dm pipeline.DependencyManager) error {
			var store storage.Store
			var storageConfig dto.ProjectStorageConfig

			err := dm.Invoke(func(awsConfigDep dto.ProjectStorageConfig, storeDep storage.Store) {
				store = storeDep
				storageConfig = awsConfigDep
			})
			if err != nil {
				return err
			}
			projectStorageKey := storageConfig.ProjectStorageKey
			projectStorageBucket := storageConfig.ProjectStorageBucket

			projectTarFile, err := store.Get(context.Background(), projectStorageBucket, projectStorageKey)
			if err != nil {
				return err
			}
			dm.AddDependency(func() dto.ProjectFileData {
				return dto.ProjectFileData{ProjectTarFile: projectTarFile, DirName: "app"}
			})
			slog.Debug("Files pulled from storage successfully")
			return nil
		},
	}
}
