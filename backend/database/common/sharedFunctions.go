package common

import "github.com/xomatix/silly-syntax-backend-bonanza/database"

func ExecuteNonQuery(q string) (int64, error) {
	return database.ExecuteNonQuery(q)
}

func ExecuteQuery(q string) ([]map[string]interface{}, error) {
	return database.ExecuteQuery(q)
}
