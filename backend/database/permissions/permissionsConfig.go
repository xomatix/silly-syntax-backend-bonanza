package permissions

import (
	"encoding/json"
	"fmt"

	"github.com/xomatix/silly-syntax-backend-bonanza/database"
	querygenerators "github.com/xomatix/silly-syntax-backend-bonanza/database/queryGenerators"
)

type TablePermission struct {
	TableName string `json:"tableName"`
	Create    string `json:"c"`
	Read      string `json:"r"`
	Update    string `json:"u"`
	Delete    string `json:"d"`
}

var (
	tablesPermissionsConfig map[string]TablePermission
)

func GetTablePermissions(tableName string) (TablePermission, error) {
	if len(tablesPermissionsConfig) == 0 {
		LoadTablesPermissions()
	}
	if conf, ok := tablesPermissionsConfig[tableName]; ok {
		return conf, nil
	}
	return TablePermission{}, fmt.Errorf("permission for table %s does not exist", tableName)
}

func GetAllTablePermissions() map[string]TablePermission {
	if len(tablesPermissionsConfig) == 0 {
		if err := LoadTablesPermissions(); err != nil {
			fmt.Println(err)

		}
	}

	return tablesPermissionsConfig
}

func LoadTablesPermissions() error {
	selectQuery, err := querygenerators.SelectQueryCreator(querygenerators.SelectQueryCreator{CollectionName: "tables_permissions"}).GetQuery()
	if err != nil {
		return fmt.Errorf("failed to load tables permissions: %v 1", err)
	}

	res, err := database.ExecuteQuery(selectQuery)
	if err != nil {
		return fmt.Errorf("failed to load tables permissions: %v 2", err)
	}

	loadedPermissions := make(map[string]TablePermission)
	for _, individualRow := range res {
		tp := TablePermission{}
		convertedVal, err := json.Marshal(individualRow)
		if err != nil {
			return fmt.Errorf("failed to load tables permissions: %v 3", err)
		}
		json.Unmarshal(convertedVal, &tp)
		loadedPermissions[tp.TableName] = tp
	}

	tablesPermissionsConfig = loadedPermissions
	return nil
}
