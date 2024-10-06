package pluginfunctions

import (
	"fmt"
	"strings"
)

type PluginLoader struct {
	Plugins  []string
	Triggers map[string]map[string][](func(map[string]interface{}, *map[string]interface{}) error)
	Api      map[string](func(map[string]interface{}, int64) (map[string]interface{}, error))
}

var instance *PluginLoader

func GetPluginLoader() *PluginLoader {
	if instance == nil {
		instance = &PluginLoader{
			Plugins:  make([]string, 0),
			Triggers: make(map[string]map[string][](func(map[string]interface{}, *map[string]interface{}) error)),
			Api:      make(map[string](func(map[string]interface{}, int64) (map[string]interface{}, error))),
		}
	}
	return instance
}

func (pl *PluginLoader) AddTrigger(caller string, f []func(map[string]interface{}, *map[string]interface{}) error) {
	call := strings.Split(caller, "/")
	if _, ok := pl.Triggers[call[0]]; !ok {
		pl.Triggers[call[0]] = make(map[string][](func(map[string]interface{}, *map[string]interface{}) error))
	}
	if _, ok := pl.Triggers[call[0]][call[1]]; !ok {
		pl.Triggers[call[0]][call[1]] = make([](func(map[string]interface{}, *map[string]interface{}) error), 0)
	}

	pl.Triggers[call[0]][call[1]] = append(pl.Triggers[call[0]][call[1]], f...)
}

func (pl *PluginLoader) AddApiRoute(caller string, f func(map[string]interface{}, int64) (map[string]interface{}, error)) {
	if _, ok := pl.Api[caller]; ok {
		fmt.Printf("Api route %s already exists overwriting\n", caller)
		return
	}
	pl.Api[caller] = f
}
