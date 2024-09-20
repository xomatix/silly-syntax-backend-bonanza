package querygenerators

import (
	"fmt"
	"strings"
)

type SelectQueryCreator struct {
	CollectionName string  `json:"collectionName"`
	ID             []int64 `json:"id"`
	Limit          int     `json:"limit"`
	Size           int     `json:"size"`
	Filter         string  `json:"filter"`
}

func (q SelectQueryCreator) GetQuery() (string, error) {

	where := ""

	countIds := len(q.ID)
	if len(q.ID) > 0 {
		q.Limit = countIds

		idsArr := make([]string, countIds)

		for i := 0; i < countIds; i++ {
			idsArr[i] = fmt.Sprintf("%d", q.ID[i])
		}

		where = fmt.Sprintf("id IN (%s)", strings.Join(idsArr, ","))
	}

	if len(q.Filter) > 0 {
		filters := ""

		q.Filter = strings.ReplaceAll(q.Filter, "&&", "AND")
		q.Filter = strings.ReplaceAll(q.Filter, "||", "OR")
		q.Filter = strings.ReplaceAll(q.Filter, "~", "LIKE")

		filters = q.Filter

		if len(where) > 0 {
			filters = fmt.Sprintf(" AND (%s)", filters)
		}
		where = fmt.Sprintf("%s %s", where, filters)
	}

	if len(where) > 0 {
		where = fmt.Sprintf(" WHERE %s", where)
	}

	limit := ""
	if q.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", q.Limit)
	}

	size := ""
	if q.Size > 0 && q.Limit > 0 {
		size = fmt.Sprintf(" OFFSET %d", q.Size)
	}

	query := fmt.Sprintf("SELECT * FROM %s %s ORDER BY created DESC %s %s;", q.CollectionName, where, limit, size)
	return query, nil
}
