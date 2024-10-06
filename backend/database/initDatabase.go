package database

import (
	"fmt"

	"github.com/xomatix/silly-syntax-backend-bonanza/database/authentication"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/database_config"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/database_functions"
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
	_, err := database_functions.ExecuteNonQuery(qInitTablesConfig)
	if err != nil {
		return err
	}

	return nil
}

func InitDatabasePermissions() error {

	tpName := "tables_permissions"
	//tables_permissions
	if _, err := database_config.GetTableConfig(tpName); err != nil {
		err := CrateTable(tpName)
		if err != nil {
			return err
		}
		AddColumnToTable(tpName, database_config.ColumnConfig{
			Name:     "tableName",
			DataType: database_config.DTTEXT,
			Unique:   true,
			NotNull:  true,
		})
		AddColumnToTable(tpName, database_config.ColumnConfig{
			Name:     "c",
			DataType: database_config.DTTEXT,
		})
		AddColumnToTable(tpName, database_config.ColumnConfig{
			Name:     "r",
			DataType: database_config.DTTEXT,
		})
		AddColumnToTable(tpName, database_config.ColumnConfig{
			Name:     "u",
			DataType: database_config.DTTEXT,
		})
		AddColumnToTable(tpName, database_config.ColumnConfig{
			Name:     "d",
			DataType: database_config.DTTEXT,
		})
		qPermissionsInit := fmt.Sprintf(`INSERT INTO tables_permissions (tableName, c,r,u,d) VALUES ('%s', '', '', '', '');`, tpName)
		_, err = database_functions.ExecuteNonQuery(qPermissionsInit)
		if err != nil {
			return err
		}
	}

	//users
	if _, err := database_config.GetTableConfig("users"); err != nil {
		err := CrateTable("users")
		if err != nil {
			return err
		}
		AddColumnToTable("users", database_config.ColumnConfig{
			Name:     "password",
			DataType: database_config.DTTEXT,
			NotNull:  true,
			Unique:   true,
		})
		AddColumnToTable("users", database_config.ColumnConfig{
			Name:     "username",
			DataType: database_config.DTTEXT,
			NotNull:  true,
		})
		AddColumnToTable("users", database_config.ColumnConfig{
			Name:     "email",
			DataType: database_config.DTTEXT,
			Unique:   true,
		})

		hPassword, _ := authentication.HashPassword("admin")
		qInsertAdmin := fmt.Sprintf(`INSERT INTO users (username, email, password) VALUES ('admin', 'admin', '%s');`, hPassword)

		_, err = database_functions.ExecuteNonQuery(qInsertAdmin)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitDatabaseViews() error {

	tpName := "views"
	//views
	if _, err := database_config.GetTableConfig(tpName); err != nil {
		err := CrateTable(tpName)
		if err != nil {
			return err
		}
		AddColumnToTable(tpName, database_config.ColumnConfig{
			Name:     "name",
			DataType: database_config.DTTEXT,
			Unique:   true,
			NotNull:  true,
		})
		AddColumnToTable(tpName, database_config.ColumnConfig{
			Name:     "query",
			DataType: database_config.DTTEXT,
		})
	}

	return nil
}
