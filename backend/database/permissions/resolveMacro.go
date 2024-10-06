package permissions

import (
	"fmt"

	"strings"

	"github.com/xomatix/silly-syntax-backend-bonanza/database/database_config"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/database_functions"
	querygenerators "github.com/xomatix/silly-syntax-backend-bonanza/database/queryGenerators"
)

func ResolvePermissionsMacro(macro string, userID int64) string {
	q := querygenerators.SelectQueryCreator{
		CollectionName: "users",
		ID:             []int64{userID},
	}

	userConfig, _ := database_config.GetTableConfig("users")

	query, err := q.GetQuery()
	if err != nil {
		return macro
	}

	result, err := database_functions.ExecuteQuery(query)
	if err != nil {
		return macro
	}

	if len(result) == 0 {
		return macro
	}

	for key, value := range result[0] {
		resolved := false
		for _, v := range userConfig.Columns {
			if v.Name == key && (v.DataType == database_config.DTTEXT || v.DataType == database_config.DTDATETIME) {
				macro = strings.ReplaceAll(macro, fmt.Sprintf("@user.%s", key), fmt.Sprintf("'%v'", value))
				resolved = true
				break
			}
		}
		if !resolved {
			macro = strings.ReplaceAll(macro, fmt.Sprintf("@user.%s", key), fmt.Sprintf("%v", value))
		}
	}

	return macro
}
