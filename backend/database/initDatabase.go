package database

import (
	"fmt"

	"github.com/xomatix/silly-syntax-backend-bonanza/database/authentication"
)

func InitDatabase() error {

	//check if database exists
	// res, err := ExecuteQuery("SELECT count(*)>0 as exists FROM settings WHERE key='database_version' limit 1;")

	// if err == nil && len(res) > 0 && res[0]["exists"].(int64) > 0 {
	// 	return nil
	// }

	// qInitSettings := `CREATE TABLE IF NOT EXISTS settings (
	// 	id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 	key TEXT,
	// 	value TEXT
	// 	);`
	// err = ExecuteNonQuery(qInitSettings)
	// if err != nil {
	// 	return err
	// }
	// qInitSettingsIndex := `CREATE INDEX IF NOT EXISTS settings_key_idx ON settings (key);`
	// err = ExecuteNonQuery(qInitSettingsIndex)
	// if err != nil {
	// 	return err
	// }

	qInitTablesConfig := `CREATE TABLE IF NOT EXISTS tables_config (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT,
		config TEXT
	);`
	_, err := ExecuteNonQuery(qInitTablesConfig)
	if err != nil {
		return err
	}

	return nil
}

func InitDatabasePermissions() error {

	tpName := "tables_permissions"
	//tables_permissions
	if _, err := GetTableConfig(tpName); err != nil {
		err := CrateTable(tpName)
		if err != nil {
			return err
		}
		AddColumnToTable(tpName, ColumnConfig{
			Name:     "tableName",
			DataType: DTTEXT,
			Unique:   true,
			NotNull:  true,
		})
		AddColumnToTable(tpName, ColumnConfig{
			Name:     "c",
			DataType: DTTEXT,
		})
		AddColumnToTable(tpName, ColumnConfig{
			Name:     "r",
			DataType: DTTEXT,
		})
		AddColumnToTable(tpName, ColumnConfig{
			Name:     "u",
			DataType: DTTEXT,
		})
		AddColumnToTable(tpName, ColumnConfig{
			Name:     "d",
			DataType: DTTEXT,
		})
		qPermissionsInit := fmt.Sprintf(`INSERT INTO tables_permissions (tableName, c,r,u,d) VALUES ('%s', '', '', '', '');`, tpName)
		_, err = ExecuteNonQuery(qPermissionsInit)
		if err != nil {
			return err
		}
	}

	//users
	if _, err := GetTableConfig("users"); err != nil {
		err := CrateTable("users")
		if err != nil {
			return err
		}
		AddColumnToTable("users", ColumnConfig{
			Name:     "password",
			DataType: DTTEXT,
			NotNull:  true,
			Unique:   true,
		})
		AddColumnToTable("users", ColumnConfig{
			Name:     "username",
			DataType: DTTEXT,
			NotNull:  true,
		})
		AddColumnToTable("users", ColumnConfig{
			Name:     "email",
			DataType: DTTEXT,
			Unique:   true,
		})

		hPassword, _ := authentication.HashPassword("admin")
		qInsertAdmin := fmt.Sprintf(`INSERT INTO users (username, email, password) VALUES ('admin', 'admin', '%s');`, hPassword)

		_, err = ExecuteNonQuery(qInsertAdmin)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitDatabaseViews() error {

	tpName := "views"
	//views
	if _, err := GetTableConfig(tpName); err != nil {
		err := CrateTable(tpName)
		if err != nil {
			return err
		}
		AddColumnToTable(tpName, ColumnConfig{
			Name:     "name",
			DataType: DTTEXT,
			Unique:   true,
			NotNull:  true,
		})
		AddColumnToTable(tpName, ColumnConfig{
			Name:     "query",
			DataType: DTTEXT,
		})
	}

	return nil
}
