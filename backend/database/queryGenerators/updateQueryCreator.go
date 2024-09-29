package querygenerators

import (
	"fmt"
	"strings"

	"github.com/xomatix/silly-syntax-backend-bonanza/database"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/authentication"
)

type UpdateQueryCreator struct {
	CollectionName string         `json:"collectionName"`
	ID             int64          `json:"id"`
	Values         map[string]any `json:"values"`
	Filter         string
}

func (q UpdateQueryCreator) SelectQuery() (string, error) {
	selectQueryCreator := SelectQueryCreator{
		CollectionName: q.CollectionName,
		ID:             []int64{q.ID},
	}
	return selectQueryCreator.GetQuery()
}

func (q UpdateQueryCreator) UpdateQuery() (string, error) {

	tableConfig, err := database.GetTableConfig(q.CollectionName)
	if err != nil {
		return "", err
	}

	var sqlValues []string

	for _, column := range tableConfig.Columns {
		if column.Name == "id" || column.Name == "created" || column.Name == "updated" {
			continue
		}
		val, exists := q.Values[column.Name]
		if !exists || val == nil {
			continue
		}
		if column.Name != "password" && column.NotNull && (len(fmt.Sprintf("%v", val)) == 0) {
			return "", fmt.Errorf("column %s cannot be null", column.Name)
		}
		if column.DataType == database.DTDOUBLE || column.DataType == database.DTINTEGER || column.DataType == database.DTBOOLEAN {
			sqlValues = append(sqlValues, fmt.Sprintf("%s = %v", column.Name, val))
		} else if column.DataType == database.DTTEXT && !(len(val.(string)) == 0) && q.CollectionName == "users" && column.Name == "password" {
			hashedPassword, err := authentication.HashPassword(val.(string))
			if err != nil {
				return "", fmt.Errorf("failed to hash password: %v", err)
			}
			sqlValues = append(sqlValues, fmt.Sprintf("'%s'", hashedPassword))
		} else if column.DataType == database.DTREFERENCE && len(val.(string)) > 0 {
			if int(val.(float64)) != 0 || column.NotNull {
				isNotPresent := checkIfForeignKeysExist(column.ReferenceTable, fmt.Sprintf("%d", int(val.(float64))))
				if isNotPresent != nil {
					return "", isNotPresent
				}
			}
			sqlValues = append(sqlValues, fmt.Sprintf("%s = %d", column.Name, int(val.(float64))))
		} else {
			val = strings.ReplaceAll(val.(string), "'", "''")
			sqlValues = append(sqlValues, fmt.Sprintf("%s = '%s'", column.Name, val))
		}
	}

	joinedValues := strings.Join(sqlValues, ",")

	if q.Filter != "" {
		q.Filter = fmt.Sprintf("AND (%s)", q.Filter)
	}
	return fmt.Sprintf("UPDATE %s SET %s WHERE (id = %d) %s;", q.CollectionName, joinedValues, q.ID, q.Filter), nil
}

// checkIfForeignKeysExist checks if a foreign key exists in a given collection.
//
// collectionName is the name of the collection to check, and id is the id of the foreign key.
// Returns an error if the foreign key does not exist, otherwise returns nil.
func checkIfForeignKeysExist(collectionName string, id string) error {
	q := fmt.Sprintf("SELECT count(*) > 0 as e FROM %s WHERE id = %s LIMIT 1;", collectionName, id)

	res, err := database.ExecuteQuery(q)

	if err == nil && len(res) > 0 && res[0]["e"].(int64) > 0 {
		return nil
	}
	return fmt.Errorf("record with id %s does not exist in collection '%s'", id, collectionName)
}
