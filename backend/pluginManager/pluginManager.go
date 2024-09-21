package pluginmanager

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/xomatix/silly-syntax-backend-bonanza/database"
)

// loads all plugins
type PluginLoader struct {
	Plugins  []string
	Triggers map[string]map[string][](func(map[string]interface{}, *map[string]interface{}) error)
}

var instance *PluginLoader

func GetPluginLoader() *PluginLoader {
	if instance == nil {
		instance = &PluginLoader{
			Plugins:  make([]string, 0),
			Triggers: make(map[string]map[string][](func(map[string]interface{}, *map[string]interface{}) error)),
		}
	}
	return instance
}

func (pl *PluginLoader) AddTrigger(caller string, f []func(map[string]interface{}, *map[string]interface{}) error) {
	call := strings.Split(caller, "/")
	if _, ok := pl.Triggers[call[0]]; ok == false {
		pl.Triggers[call[0]] = make(map[string][](func(map[string]interface{}, *map[string]interface{}) error))
	}
	if _, ok := pl.Triggers[call[0]][call[1]]; ok == false {
		pl.Triggers[call[0]][call[1]] = make([](func(map[string]interface{}, *map[string]interface{}) error), 0)
	}

	pl.Triggers[call[0]][call[1]] = append(pl.Triggers[call[0]][call[1]], f...)
}

func (pl *PluginLoader) LoadPlugins(mux *http.ServeMux) error {
	pluginDir := "./plugins"
	files, err := os.ReadDir(pluginDir)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", pluginDir, err)
	}

	executeQuery := func(q string) ([]map[string]interface{}, error) {
		return database.ExecuteQuery(q)
	}

	executeNonQuery := func(q string) (int64, error) {
		return database.ExecuteNonQuery(q)
	}

	// Create a map to pass as context
	plArgs := map[string]interface{}{
		"Mux":             mux,
		"ExecuteQuery":    executeQuery,
		"ExecuteNonQuery": executeNonQuery,
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".so" {
			pluginPath := filepath.Join(pluginDir, file.Name())
			err = pl.loadPlugin(pluginPath, plArgs)
			if err != nil {
				log.Printf("Failed to load plugin %s: %v", pluginPath, err)
			} else {
				log.Printf("Successfully loaded plugin %s", pluginPath)
			}
		}
	}
	return nil
}

func (pl *PluginLoader) loadPlugin(pluginPath string, ctx map[string]interface{}) error {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %v", err)
	}

	initPluginSymbol, err := p.Lookup("InitPlugin")
	if err != nil {
		return fmt.Errorf("failed to lookup InitPlugin symbol: %v", err)
	}

	initPluginFunc, ok := initPluginSymbol.(func(map[string]interface{}) (map[string][](func(map[string]interface{}, *map[string]interface{}) error), error))
	if !ok {
		return fmt.Errorf("plugin does not have the expected InitPlugin function signature")
	}

	triggers, err := initPluginFunc(ctx)

	for k, f := range triggers {
		pl.AddTrigger(k, f)
	}

	pluginName := filepath.Base(pluginPath)
	pl.Plugins = append(pl.Plugins, pluginName)

	return nil
}
