// Package contentregistry is the Content Registry BASE plugin for togo: One declarative map of every content type (table, slug col, translatable fields, feed-kind, icon, section, searchable) that feed / board / seo / mcp / search / compose all read — replacing five duplicated hardcoded maps.
//
// Driver plugins self-register a Descriptor via contentregistry.Register in their init();
// an app installs a driver by blank-importing it. This base owns the registry and
// exposes it at GET /api/registry.
package contentregistry

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"

	"github.com/togo-framework/togo"
)

// Descriptor is one registered Content Registry item contributed by a driver.
type Descriptor struct {
	Key         string `json:"key"`
	Title       string `json:"title"`
	DefaultZone string `json:"defaultZone,omitempty"`
	Asset       string `json:"asset,omitempty"` // URL of the driver's built web asset
}

var (
	mu       sync.RWMutex
	registry = map[string]Descriptor{}
)

// Register installs a descriptor. Drivers call this from their init().
func Register(d Descriptor) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := registry[d.Key]; exists {
		return
	}
	registry[d.Key] = d
}

// All returns every registered descriptor, sorted by key.
func All() []Descriptor {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]Descriptor, 0, len(registry))
	for _, d := range registry {
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

func init() {
	togo.RegisterProviderFunc("content-registry", togo.PriorityService, func(k *togo.Kernel) error {
		k.Router.Get("/api/registry", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"registry": All()})
		})
		return nil
	})
}
