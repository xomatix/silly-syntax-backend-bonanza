package querygenerators

import (
	"fmt"

	"github.com/xomatix/silly-syntax-backend-bonanza/database"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/authentication"

	"strings"
)

type InsertQueryCreator struct {
	CollectionName string         `json:"collectionName"`
	Values         map[string]any `json:"values"`
	Filter         string
}

// InsertQuery generates an SQL insert query based on the values provided in the InsertQueryCreator instance.
//
// It retrieves the table configuration for the specified collection, constructs the SQL query based on the column values,
// and returns the generated insert query string along with a potential error and guid of the inserted row.
func (q InsertQueryCreator) InsertQuery() (string, error) {

	tableConfig, err := database.GetTableConfig(q.CollectionName)
	if err != nil {
		return "", err
	}

	var sqlValues []string
	var sqlColumns []string

	for _, column := range tableConfig.Columns {
		initialLen := len(sqlValues)
		if column.Name == "id" || column.Name == "created" || column.Name == "updated" {
			continue
		}
		val, exists := q.Values[column.Name]
		if !exists {
			continue
		}
		if column.DataType == database.DTDOUBLE || column.DataType == database.DTINTEGER || column.DataType == database.DTBOOLEAN || column.DataType == database.DTREFERENCE {
			if len(fmt.Sprintf("%v", val)) > 0 {
				sqlValues = append(sqlValues, fmt.Sprintf("%v", val))
			}
		} else if column.DataType == database.DTTEXT && q.CollectionName == "users" && column.Name == "password" {
			hashedPassword, err := authentication.HashPassword(val.(string))
			if err != nil {
				return "", fmt.Errorf("failed to hash password: %v", err)
			}
			sqlValues = append(sqlValues, fmt.Sprintf("'%s'", hashedPassword))
		} else {
			val = strings.ReplaceAll(val.(string), "'", "''")
			sqlValues = append(sqlValues, fmt.Sprintf("'%s'", val))
		}
		if initialLen < len(sqlValues) {
			sqlColumns = append(sqlColumns, column.Name)
		}
	}

	joinedValues := strings.Join(sqlValues, ",")
	joinedColumns := strings.Join(sqlColumns, ",")

	if q.Filter != "" {
		q.Filter = fmt.Sprintf("WHERE %s", q.Filter)
	}

	return fmt.Sprintf("INSERT INTO %s (%s) SELECT %s %s;", q.CollectionName, joinedColumns, joinedValues, q.Filter), nil
}
