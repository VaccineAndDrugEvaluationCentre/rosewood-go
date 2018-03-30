// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package rosewood

import (
	"fmt"
	"sort"
	"sync"

	"./types"
)

type Config struct {
	Name     string
	Renderer types.Renderer
}

var (
	renderersMu sync.RWMutex
	renderers   = make(map[string]types.Renderer)
)

// Register makes an output renderer available by the provided name.
// If Register is called twice with the same name or if renderer is nil,
// it panics.
func Register(config *Config) {
	renderersMu.Lock()
	defer renderersMu.Unlock()
	if config.Renderer == nil {
		panic("rosewood: Register renderer is nil")
	}
	if _, dup := renderers[config.Name]; dup {
		panic("rosewood: Register called twice for renderer " + config.Name)
	}
	renderers[config.Name] = config.Renderer
}

func unregisterAllRenderers() {
	renderersMu.Lock()
	defer renderersMu.Unlock()
	renderers = make(map[string]types.Renderer)
}

// Renderers returns a sorted list of the names of the registered renderers.
func Renderers() []string {
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
func GetRendererByName(name string) (types.Renderer, error) {
	renderersMu.RLock()
	renderer, ok := renderers[name]
	renderersMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("rosewood: unknown renderer %q (forgotten import?)", name)
	}
	return renderer, nil
}
