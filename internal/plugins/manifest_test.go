package plugins

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Brook-sys/auxitalk/pkg/types"
)

func writeManifest(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "plugin.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	return path
}

func validManifestJSON(id string) string {
	return `{
		"id": "` + id + `",
		"name": "Mock Input",
		"version": "0.1.0",
		"runtime": "node",
		"entry": "index.js",
		"kind": "input",
		"permissions": ["event.emit"],
		"capabilities": [{"name": "conversation.observe"}]
	}`
}

func TestLoadManifest(t *testing.T) {
	path := writeManifest(t, validManifestJSON("mock-input"))

	file, err := LoadManifest(path)
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}

	if file.Manifest.ID != "mock-input" {
		t.Fatalf("expected mock-input, got %q", file.Manifest.ID)
	}
	if file.Path == "" {
		t.Fatal("expected absolute manifest path")
	}
	if file.Dir == "" {
		t.Fatal("expected manifest dir")
	}
}

func TestLoadManifestInvalidJSON(t *testing.T) {
	path := writeManifest(t, `{`)

	if _, err := LoadManifest(path); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestLoadManifestValidationError(t *testing.T) {
	path := writeManifest(t, `{"id":"missing-fields"}`)

	if _, err := LoadManifest(path); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestRegistryRegisterGetListCount(t *testing.T) {
	registry := NewRegistry()
	plugin := Plugin{
		ManifestPath: "/tmp/plugin.json",
		RootDir:      "/tmp",
		Manifest: types.PluginManifest{
			ID:      "mock-input",
			Name:    "Mock Input",
			Version: "0.1.0",
			Runtime: "node",
			Entry:   "index.js",
			Kind:    types.PluginKindInput,
		},
	}

	if err := registry.Register(plugin); err != nil {
		t.Fatalf("register: %v", err)
	}

	found, ok := registry.Get("mock-input")
	if !ok {
		t.Fatal("expected plugin to exist")
	}
	if found.Manifest.ID != "mock-input" {
		t.Fatalf("unexpected plugin id %q", found.Manifest.ID)
	}
	if registry.Count() != 1 {
		t.Fatalf("expected count 1, got %d", registry.Count())
	}
	if len(registry.List()) != 1 {
		t.Fatal("expected one listed plugin")
	}
}

func TestRegistryRejectsDuplicateID(t *testing.T) {
	registry := NewRegistry()
	manifest := types.PluginManifest{
		ID:      "mock-input",
		Name:    "Mock Input",
		Version: "0.1.0",
		Runtime: "node",
		Entry:   "index.js",
		Kind:    types.PluginKindInput,
	}

	if err := registry.Register(Plugin{Manifest: manifest}); err != nil {
		t.Fatalf("first register: %v", err)
	}

	err := registry.Register(Plugin{Manifest: manifest})
	if !errors.Is(err, ErrPluginAlreadyRegistered) {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestRegistryRegisterManifest(t *testing.T) {
	path := writeManifest(t, validManifestJSON("mock-input"))
	file, err := LoadManifest(path)
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}

	registry := NewRegistry()
	if err := registry.RegisterManifest(file); err != nil {
		t.Fatalf("register manifest: %v", err)
	}

	plugin, ok := registry.Get("mock-input")
	if !ok {
		t.Fatal("expected registered manifest")
	}
	if plugin.ManifestPath != file.Path {
		t.Fatalf("expected manifest path %q, got %q", file.Path, plugin.ManifestPath)
	}
}
