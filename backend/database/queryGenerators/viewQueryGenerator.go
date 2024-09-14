package querygenerators

import (
	"fmt"
)

type ViewQueryCreator struct {
	Query string
	// ID             []int64 `json:"id"`
	Limit int `json:"limit"`
	Size  int `json:"size"`
	// Filter         string  `json:"filter"`
}

func (q ViewQueryCreator) GetViewQuery() (string, error) {

	limit := ""
	if q.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", q.Limit)
	}

	size := ""
	if q.Size > 0 && q.Limit > 0 {
		size = fmt.Sprintf(" OFFSET %d", q.Size)
	}

	query := fmt.Sprintf("%s %s %s;", q.Query, limit, size)
	return query, nil
}
