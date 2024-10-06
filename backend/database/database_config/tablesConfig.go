package database_config

import (
	"encoding/json"
	"fmt"

	"github.com/xomatix/silly-syntax-backend-bonanza/database/database_functions"
)

type TableConfig struct {
	Name    string         `json:"name"`
	Columns []ColumnConfig `json:"columns"`
}

type ColumnConfig struct {
	Name           string   `json:"name"`
	DataType       DataType `json:"dataType"`
	ReferenceTable string   `json:"refTable"`
	NotNull        bool     `json:"notNull"`
	Unique         bool     `json:"unique"`
}

type DataType string

// Constants representing enum values
const (
	DTINTEGER   DataType = "INTEGER"
	DTDOUBLE    DataType = "DOUBLE"
	DTTEXT      DataType = "TEXT"
	DTBOOLEAN   DataType = "BOOLEAN"
	DTDATETIME  DataType = "DATETIME"
	DTREFERENCE DataType = "REFERENCE"
)

var (
	tablesConfig map[string]TableConfig
)

// GetTablesConfig retrieves the tables configuration map.
//
// No parameters.
// Returns a map of string keys to TableConfig values. To be used in query builder
func GetTablesConfig() map[string]TableConfig {
	return tablesConfig
}

func SetTablesConfig(newConfig map[string]TableConfig) {
	tablesConfig = newConfig
}

// GetTableConfig retrieves the TableConfig for a specific table by its name.
//
// tableName: the name of the table for which to retrieve the TableConfig.
// Returns the TableConfig for the specified table. To be used in query builder
func GetTableConfig(tableName string) (TableConfig, error) {
	if conf, ok := tablesConfig[tableName]; ok {
		return conf, nil
	}
	return TableConfig{}, fmt.Errorf("table %s does not exist", tableName)
}

// LoadTablesConfig retrieves the table configuration from the database and loads it into the tablesConfig map.
// Call after table modifications
//
// No parameters.
// No return values.
func LoadTablesConfig() {
	res, err := database_functions.ExecuteQuery("SELECT id, config FROM tables_config;")
	if err != nil {
		fmt.Printf("failed to load tables from database: %v", err)
	}

	tablesConfig = make(map[string]TableConfig)
	for _, row := range res {
		tabConf, err := jsonToTableConfig(row["config"].(string))
		if err != nil {
			fmt.Printf("failed to load table config: %v", err)
		}
		tablesConfig[tabConf.Name] = tabConf
	}
}

// jsonToTableConfig unmarshalls JSON string to TableConfig struct.
//
// jsonStr: string containing JSON data to unmarshal.
// Returns TableConfig struct and error.
func jsonToTableConfig(jsonStr string) (TableConfig, error) {
	var tabConf TableConfig
	err := json.Unmarshal([]byte(jsonStr), &tabConf)
	if err != nil {
		return TableConfig{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	return tabConf, nil
}
