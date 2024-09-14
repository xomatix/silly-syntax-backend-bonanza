package permissions

import (
	"fmt"
	"silly-syntax-backend-bonanza/database"
	querygenerators "silly-syntax-backend-bonanza/database/queryGenerators"
	"strings"
)

func ResolvePermissionsMacro(macro string, userID int64) string {
	q := querygenerators.SelectQueryCreator{
		CollectionName: "users",
		ID:             []int64{userID},
	}

	userConfig, _ := database.GetTableConfig("users")

	query, err := q.GetQuery()
	if err != nil {
		return macro
	}

	result, err := database.ExecuteQuery(query)
	if err != nil {
		return macro
	}

	if len(result) == 0 {
		return macro
	}

	for key, value := range result[0] {
		resolved := false
		for _, v := range userConfig.Columns {
			if v.Name == key && (v.DataType == database.DTTEXT || v.DataType == database.DTDATETIME) {
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
