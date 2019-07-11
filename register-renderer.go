// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package rosewood

import (
	"fmt"
	"sort"
	"sync"

	"github.com/drgo/rosewood/table"
)

//RendererFactory function signature for a func that creates a renderer
type RendererFactory func() (table.Renderer, error)

//RendererConfig holds info on this renderer
type RendererConfig struct {
	Name     string
	Renderer RendererFactory
}

var (
	renderersMu sync.RWMutex
	renderers   = make(map[string]*RendererConfig)
)

// RegisterRenderer makes an output renderer available by the provided name.
// If Register is called twice with the same name or if renderer is nil,
// it panics.
func RegisterRenderer(config *RendererConfig) {
	renderersMu.Lock()
	defer renderersMu.Unlock()
	if config.Renderer == nil {
		panic(fmt.Sprintf("rosewood: failed to register renderer %s: renderer is nil", config.Name))
	}
	if _, dup := renderers[config.Name]; dup {
		panic("rosewood: Register called twice for renderer " + config.Name)
	}
	renderers[config.Name] = config
}

func unregisterAllRenderers() {
	renderersMu.Lock()
	defer renderersMu.Unlock()
	renderers = make(map[string]*RendererConfig)
}

// GetRenderersList returns a sorted list of the names of the registered renderers.
func GetRenderersList() []string {
	renderersMu.RLock()
	defer renderersMu.RUnlock()
	var list []string
	for name := range renderers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

// GetRendererByName returns a renderer specified by its name
func GetRendererByName(name string) (table.Renderer, error) {
	renderersMu.RLock()
	rendererConfig, ok := renderers[name]
	renderersMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("rosewood: unknown renderer %q (forgotten import?)", name)
	}
	return rendererConfig.Renderer()
}
