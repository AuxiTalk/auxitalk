package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Brook-sys/auxitalk/pkg/types"
)

type ManifestFile struct {
	Path     string
	Dir      string
	Manifest types.PluginManifest
}

func LoadManifest(path string) (ManifestFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ManifestFile{}, fmt.Errorf("read plugin manifest: %w", err)
	}

	var manifest types.PluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return ManifestFile{}, fmt.Errorf("parse plugin manifest: %w", err)
	}

	if err := manifest.Validate(); err != nil {
		return ManifestFile{}, fmt.Errorf("validate plugin manifest: %w", err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return ManifestFile{}, fmt.Errorf("resolve plugin manifest path: %w", err)
	}

	return ManifestFile{
		Path:     absPath,
		Dir:      filepath.Dir(absPath),
		Manifest: manifest,
	}, nil
}
