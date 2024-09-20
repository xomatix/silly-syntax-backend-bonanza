package pluginHandler

import (
	"log"
	"os"
	"path/filepath"
	"plugin"
)

// todo
// 1. Add support for api routes
// 2. Add support for triggers
//  1. after add delete update
//
// 3. Add support for exported functions
// 4. Add support for bot tasks
func LoadPlugins() map[string](func(*map[string]interface{}) error) {
	pluginDir := "./plugins"
	plugins := make(map[string](func(*map[string]interface{}) error))

	err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v\n", path, err)
			return nil
		}

		if filepath.Ext(path) == ".so" {
			p, err := plugin.Open(path)
			if err != nil {
				log.Printf("Error opening plugin %q: %v\n", path, err)
				return nil
			}

			// Example function signature: func MyFunc(*map[string]interface{}) error
			symbolName := "MyFunc"
			symbol, err := p.Lookup(symbolName)
			if err != nil {
				log.Printf("Error looking up symbol %q in plugin %q: %v\n", symbolName, path, err)
				return nil
			}

			// Cast symbol to the correct function signature
			function, ok := symbol.(func(*map[string]interface{}) error)
			if !ok {
				log.Printf("Symbol %q from plugin %q has an incompatible type\n", symbolName, path)
				return nil
			}

			// Store the function in the map
			plugins[path] = function
			log.Printf("Successfully loaded function %q from plugin %q\n", symbolName, path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the plugin directory: %v", err)
	}

	return plugins
}
