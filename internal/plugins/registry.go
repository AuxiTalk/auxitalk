package plugins

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/Brook-sys/auxitalk/pkg/types"
)

var ErrPluginAlreadyRegistered = errors.New("plugin already registered")

type Plugin struct {
	ManifestPath string
	RootDir      string
	Manifest     types.PluginManifest
}

type Registry struct {
	mu      sync.RWMutex
	plugins map[string]Plugin
}

func NewRegistry() *Registry {
	return &Registry{plugins: make(map[string]Plugin)}
}

func (r *Registry) Register(plugin Plugin) error {
	if err := plugin.Manifest.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	id := plugin.Manifest.ID
	if _, exists := r.plugins[id]; exists {
		return fmt.Errorf("%w: %s", ErrPluginAlreadyRegistered, id)
	}

	r.plugins[id] = plugin
	return nil
}

func (r *Registry) RegisterManifest(file ManifestFile) error {
	return r.Register(Plugin{
		ManifestPath: file.Path,
		RootDir:      file.Dir,
		Manifest:     file.Manifest,
	})
}

func (r *Registry) Get(id string) (Plugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, ok := r.plugins[id]
	return plugin, ok
}

func (r *Registry) List() []Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make([]Plugin, 0, len(r.plugins))
	for _, plugin := range r.plugins {
		plugins = append(plugins, plugin)
	}

	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Manifest.ID < plugins[j].Manifest.ID
	})

	return plugins
}

func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.plugins)
}
