package pluginmanager

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plugin"

	"github.com/xomatix/silly-syntax-backend-bonanza/api/collection"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/authentication"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/database_functions"
	querygenerators "github.com/xomatix/silly-syntax-backend-bonanza/database/queryGenerators"
	pluginfunctions "github.com/xomatix/silly-syntax-backend-bonanza/pluginManager/plugin_functions"
)

// loads all plugins

func LoadPlugins(mux *http.ServeMux) error {
	pluginDir := "./plugins"
	files, err := os.ReadDir(pluginDir)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", pluginDir, err)
	}

	executeQuery := func(q string) ([]map[string]interface{}, error) {
		return database_functions.ExecuteQuery(q)
	}

	executeNonQuery := func(q string) (int64, error) {
		return database_functions.ExecuteNonQuery(q)
	}

	resolveCookiesWithUserFromHeader := func(r *http.Request) (int64, string, error) {
		return authentication.ResolveCookiesWithUserFromHeader(r)
	}

	collectionPostRaw := func(queryRaw map[string]interface{}, userID int64) map[string]interface{} {
		fmt.Printf("queryRaw: %v\n", queryRaw)
		query := querygenerators.InsertQueryCreator{
			CollectionName: queryRaw["collectionName"].(string),
			Values:         queryRaw["values"].(map[string]any),
		}
		fmt.Printf("query: %v\n", query)

		resp := collection.CollectionPost(query, userID)

		fmt.Printf("resp im core: %v\n", resp)
		return map[string]interface{}{
			"success": resp.Success,
			"message": resp.Message,
			"data":    resp.Data,
		}
	}

	// Create a map to pass as context
	plArgs := map[string]interface{}{
		"Mux":                              mux,
		"ExecuteQuery":                     executeQuery,
		"ExecuteNonQuery":                  executeNonQuery,
		"CollectionPostRaw":                collectionPostRaw,
		"ResolveCookiesWithUserFromHeader": resolveCookiesWithUserFromHeader,
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".so" {
			pluginPath := filepath.Join(pluginDir, file.Name())
			err = loadPlugin(pluginPath, plArgs)
			if err != nil {
				log.Printf("Failed to load plugin %s: %v", pluginPath, err)
			} else {
				//log.Printf("Successfully loaded plugin %s", pluginPath)
			}
		}
	}
	return nil
}

func loadPlugin(pluginPath string, ctx map[string]interface{}) error {
	pl := pluginfunctions.GetPluginLoader()
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %v", err)
	}

	initPluginSymbol, err := p.Lookup("InitPlugin")
	if err != nil {
		return fmt.Errorf("failed to lookup InitPlugin symbol: %v", err)
	}

	initPluginFunc, ok := initPluginSymbol.(func(map[string]interface{}) (map[string][](func(map[string]interface{}, *map[string]interface{}) error), map[string](func(map[string]interface{}, int64) (map[string]interface{}, error)), error))
	if !ok {
		return fmt.Errorf("plugin does not have the expected InitPlugin function signature")
	}

	triggers, apiRoutes, err := initPluginFunc(ctx)
	pluginName := filepath.Base(pluginPath)

	if err != nil {
		return fmt.Errorf("failed to initialize plugin: %s", pluginName)
	}

	for k, f := range triggers {
		pl.AddTrigger(k, f)
	}
	for k, f := range apiRoutes {
		pl.AddApiRoute(k, f)
	}

	pl.Plugins = append(pl.Plugins, pluginName)

	return nil
}
