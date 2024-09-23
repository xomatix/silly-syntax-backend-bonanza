package querygenerators

import (
	"fmt"

	"strings"

	"github.com/xomatix/silly-syntax-backend-bonanza/database"
)

type DeleteQueryCreator struct {
	CollectionName string `json:"collectionName"`
	ID             int64  `json:"id"`
	Filter         string
}

func (q DeleteQueryCreator) DeleteQuery() (string, error) {

	_, err := database.GetTableConfig(q.CollectionName)
	if err != nil {
		return "", err
	}

	// delete all occurences of record in other tables
	referencesClear := ""

	filter := ""
	q.Filter = strings.ReplaceAll(q.Filter, "&&", "AND")
	q.Filter = strings.ReplaceAll(q.Filter, "||", "OR")
	q.Filter = strings.ReplaceAll(q.Filter, "~", "LIKE")
	if q.Filter != "" {
		filter = fmt.Sprintf(" AND (%s)", q.Filter)
	}

	tableConfig := database.GetTablesConfig()
	for _, v := range tableConfig {
		if v.Name != q.CollectionName {
			for _, column := range v.Columns {
				if column.DataType == database.DTREFERENCE && column.ReferenceTable == q.CollectionName {
					referencesClear = fmt.Sprintf("%s; UPDATE %s SET %s = NULL WHERE %s = %d %s;", referencesClear, v.Name, column.Name, column.Name, q.ID, filter)
				}
			}
		}
	}

	isNotPresent := checkIfForeignKeysExist(q.CollectionName, fmt.Sprintf("%d", q.ID))
	if isNotPresent != nil {
		return "", isNotPresent
	}

	deleteQuery := fmt.Sprintf("BEGIN;%s DELETE FROM %s WHERE (id = %d) %s; COMMIT;", referencesClear, q.CollectionName, q.ID, filter)

	return deleteQuery, nil
}
