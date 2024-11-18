package deployment

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func extractPackageJSONAsMap(tarBuffer *bytes.Buffer) (map[string]interface{}, error) {
	tarReader := tar.NewReader(tarBuffer)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(header.Name, "package.json") {
			// Read the content of package.json
			var packageJSONBuf bytes.Buffer
			if _, err := io.Copy(&packageJSONBuf, tarReader); err != nil {
				return nil, err
			}

			var packageJSONMap map[string]interface{}
			err = json.Unmarshal(packageJSONBuf.Bytes(), &packageJSONMap)
			if err != nil {
				return nil, err
			}
			return packageJSONMap, nil
		}
	}

	return nil, fmt.Errorf("package.json not found")
}

func DetectProjectType(tarBuffer *bytes.Buffer) (string, error) {
	packageJSON, err := extractPackageJSONAsMap(tarBuffer)
	if err != nil {
		return "", err
	}
	if dependencies, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
		if _, exists := dependencies["vite"]; exists {
			return "reactvitenode20", nil
		}
		if _, exists := dependencies["react-scripts"]; exists {
			return "cranode20", nil
		}
	}

	// Check devDependencies
	if devDependencies, ok := packageJSON["devDependencies"].(map[string]interface{}); ok {
		if _, exists := devDependencies["vite"]; exists {
			return "reactvitenode20", nil
		}
	}

	// Check scripts
	if scripts, ok := packageJSON["scripts"].(map[string]interface{}); ok {
		if script, exists := scripts["dev"]; exists && script == "vite" {
			return "reactvitenode20", nil
		}
		if script, exists := scripts["start"]; exists && strings.Contains(script.(string), "react-scripts") {
			return "cranode20", nil
		}
	}

	return "unknown", fmt.Errorf("unable to detect project type")
}
