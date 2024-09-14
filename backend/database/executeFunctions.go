package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "bonanza.db"

func getSQLiteConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// TODO remove this if takes time
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to verify connection: %v", err)
	}

	return db, nil
}

// ExecuteNonQuery executes a non-query SQL statement. If the statement is an INSERT, the id of the inserted row is returned.
// Otherwise, the number of affected rows is returned.
// If an error occurs, it is returned as the second return value.
// The query is executed on a connection to the local SQLite database, and the connection is closed after the query is executed.
func ExecuteNonQuery(query string) (int64, error) {
	db, err := getSQLiteConnection()
	if err != nil {
		return 0, fmt.Errorf("failed to get connection: %v", err)
	}
	result, err := db.Exec(query)
	if err != nil {
		fmt.Println(query)
	}
	defer db.Close()
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}
	var returningNumber int64
	if strings.ToLower(query[:6]) == "insert" {
		returningNumber, err = result.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("failed to execute query: %v", err)
		}
	}
	return returningNumber, nil
}

func ExecuteQuery(query string) ([]map[string]interface{}, error) {
	db, err := getSQLiteConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %v", err)
	}
	defer db.Close()
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %v", err)
	}

	var results []map[string]interface{}

	for rows.Next() {
		columnsData := make([]interface{}, len(columns))
		columnsDataPtrs := make([]interface{}, len(columns))

		for i := range columnsData {
			columnsDataPtrs[i] = &columnsData[i]
		}

		if err := rows.Scan(columnsDataPtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			rowMap[colName] = columnsData[i]
		}

		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return results, nil
}
