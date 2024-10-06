package database

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/xomatix/silly-syntax-backend-bonanza/database/database_config"
	"github.com/xomatix/silly-syntax-backend-bonanza/database/database_functions"
)

// CrateTable creates a new table in the database with the specified table name.
//
// Table by default have id text with index.
//
// created and updated timestamp
//
// tableName: the name of the table to be created.
// Returns an error if there was a problem creating the table.
func CrateTable(tableName string) error {

	actualTableConfig := database_config.GetTablesConfig()
	_, existsInConfig := actualTableConfig[tableName]
	reservedTableNames := []string{"settings", "tables_config"}

	if existsInConfig || slices.Contains(reservedTableNames, tableName) {
		return fmt.Errorf("table %s already exists", tableName)
	}

	qInitTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated DATETIME DEFAULT CURRENT_TIMESTAMP
		);`, tableName)
	_, err := database_functions.ExecuteNonQuery(qInitTable)
	if err != nil {
		return err
	}
	qInitTableIndex := fmt.Sprintf(`CREATE INDEX IF NOT EXISTS %s_idx ON %s (id);`, tableName, tableName)
	_, err = database_functions.ExecuteNonQuery(qInitTableIndex)
	if err != nil {
		return err
	}
	qUpdateTrigger := fmt.Sprintf(`CREATE TRIGGER IF NOT EXISTS trigger_%s_updated
		AFTER UPDATE ON %s
		FOR EACH ROW
		BEGIN
			UPDATE %s SET updated = CURRENT_TIMESTAMP WHERE id = OLD.id;
		END;`, tableName, tableName, tableName)
	_, err = database_functions.ExecuteNonQuery(qUpdateTrigger)
	if err != nil {
		return err
	}
	if tableName != "tables_permissions" {
		qPermissionsInit := fmt.Sprintf(`INSERT INTO tables_permissions (tableName, c,r,u,d) VALUES ('%s', '', '', '', '');`, tableName)
		_, err = database_functions.ExecuteNonQuery(qPermissionsInit)
		if err != nil {
			return err
		}
	}

	columns := []database_config.ColumnConfig{
		{
			Name:     "id",
			DataType: database_config.DTINTEGER,
			NotNull:  true,
			Unique:   true,
		},
		{
			Name:     "created",
			DataType: database_config.DTDATETIME,
			NotNull:  true,
			Unique:   true,
		},
		{
			Name:     "updated",
			DataType: database_config.DTDATETIME,
			NotNull:  true,
			Unique:   true,
		},
	}

	actualTableConfig[tableName] = database_config.TableConfig{
		Name:    tableName,
		Columns: columns,
	}
	database_config.SetTablesConfig(actualTableConfig)

	jsonStr, _ := json.Marshal(actualTableConfig[tableName])
	qInsertConfig := fmt.Sprintf(`INSERT INTO tables_config (key, config) VALUES ('%s', '%s');`, tableName, jsonStr)
	_, err = database_functions.ExecuteNonQuery(qInsertConfig)
	if err != nil {
		return err
	}

	return nil
}

// AddColumnToTable checks if the table exists, validates the data type, and adds a new column to the specified table.
//
// Parameters:
// - tableName: the name of the table to which the column will be added.
// - columnConf: the configuration of the column to be added.
// Returns an error if there was a problem adding the column.
func AddColumnToTable(tableName string, columnConf database_config.ColumnConfig) error {
	// check if table exists
	actualTableConfig, err := database_config.GetTableConfig(tableName)
	if err != nil {
		return err
	}

	// check table data type
	strDataType := ""
	switch columnConf.DataType {
	case database_config.DTBOOLEAN:
		strDataType = "BOOLEAN"
	case database_config.DTINTEGER:
		strDataType = "INTEGER"
	case database_config.DTDOUBLE:
		strDataType = "DOUBLE"
	case database_config.DTTEXT:
		strDataType = "TEXT"
	case database_config.DTDATETIME:
		strDataType = "DATETIME"
	case database_config.DTREFERENCE:
		strDataType = "INTEGER"
	}
	if strDataType == "" {
		return fmt.Errorf("invalid data type: %v", columnConf.DataType)
	}
	if columnConf.DataType == database_config.DTREFERENCE {
		_, err := database_config.GetTableConfig(columnConf.ReferenceTable)
		if err != nil {
			return fmt.Errorf("invalid reference table: %v", columnConf.ReferenceTable)
		}
	}

	if columnConf.Name == "id" || columnConf.Name == "created" || columnConf.Name == "updated" {
		return fmt.Errorf("invalid column name: %v", columnConf.Name)
	}

	for _, v := range actualTableConfig.Columns {
		if v.Name == columnConf.Name {
			return fmt.Errorf("column %s already exists in table %s", columnConf.Name, tableName)
		}
	}

	qNotNull := ""
	if columnConf.NotNull {
		qNotNull = "NOT NULL "
		if columnConf.DataType == database_config.DTREFERENCE || columnConf.DataType == database_config.DTTEXT {
			qNotNull += "DEFAULT ''"
		}
		if columnConf.DataType == database_config.DTBOOLEAN || columnConf.DataType == database_config.DTINTEGER || columnConf.DataType == database_config.DTDOUBLE || columnConf.DataType == database_config.DTDATETIME {
			qNotNull += "DEFAULT 0"
		}
	}
	qUnique := ""
	if columnConf.Unique {
		qUnique = fmt.Sprintf(`CREATE TRIGGER IF NOT EXISTS enforce_unique_%s_%s_insert
					BEFORE INSERT ON %s
					FOR EACH ROW
					WHEN NEW.%s IS NOT NULL
					BEGIN
						SELECT RAISE(ABORT, 'Duplicate value for %s: %s already exists')
						WHERE EXISTS (
							SELECT 1
							FROM %s
							WHERE %s = NEW.%s
						);
					END;`, tableName, columnConf.Name, tableName, columnConf.Name, columnConf.Name, columnConf.Name, tableName, columnConf.Name, columnConf.Name)

		qUnique += fmt.Sprintf(`CREATE TRIGGER IF NOT EXISTS enforce_unique_%s_%s_update
					BEFORE UPDATE ON %s
					FOR EACH ROW
					WHEN NEW.%s IS NOT NULL
					BEGIN
						SELECT RAISE(ABORT, 'Duplicate value for %s: %s already exists')
						WHERE EXISTS (
							SELECT 1
							FROM %s
							WHERE %s = NEW.%s
							AND id != NEW.id
						);
					END;`, tableName, columnConf.Name, tableName, columnConf.Name, columnConf.Name, columnConf.Name, tableName, columnConf.Name, columnConf.Name)
	}

	qAddColumn := fmt.Sprintf(`ALTER TABLE %s
		ADD COLUMN %s %s %s ; %s`,
		tableName,
		columnConf.Name, strDataType, qNotNull, qUnique)
	_, err = database_functions.ExecuteNonQuery(qAddColumn)
	if err != nil {
		return err
	}

	// index only mechanism for references not foreign key
	if columnConf.DataType == database_config.DTREFERENCE {
		qAddColumn := fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS idx_%s_%s ON %s(%s);`,
			columnConf.ReferenceTable, columnConf.Name, tableName, columnConf.Name)
		_, err = database_functions.ExecuteNonQuery(qAddColumn)
		if err != nil {
			return err
		}
	}

	// update table config
	actualTableConfig.Columns = append(actualTableConfig.Columns, columnConf)
	jsonStr, err := tableConfigToJson(actualTableConfig)
	if err != nil {
		return err
	}

	tablesConfig := database_config.GetTablesConfig()
	tablesConfig[tableName] = actualTableConfig
	database_config.SetTablesConfig(tablesConfig)

	qUpdateTableConfig := fmt.Sprintf(`UPDATE tables_config SET config = '%s' WHERE key = '%s';`, jsonStr, tableName)
	_, err = database_functions.ExecuteNonQuery(qUpdateTableConfig)
	if err != nil {
		return err
	}

	return nil
}

// AddColumnToTable checks if the table exists, validates the data type, and adds a new column to the specified table.
//
// Parameters:
// - tableName: the name of the table to which the column will be added.
// - columnConf: the configuration of the column to be added.
// Returns an error if there was a problem adding the column.
func RemoveColumnFromTable(tableName string, columnConfig database_config.ColumnConfig) error {
	// check if table exists
	actualTableConfig, err := database_config.GetTableConfig(tableName)
	if err != nil {
		return err
	}

	for _, v := range actualTableConfig.Columns {
		if v.Name == columnConfig.Name && v.DataType == database_config.DTREFERENCE {
			qDropIndex := fmt.Sprintf(`DROP INDEX IF EXISTS idx_%s_%s;`, v.ReferenceTable, v.Name)

			qDropUniqueConstrainInsert := fmt.Sprintf(`DROP INDEX IF EXISTS enforce_unique_%s_%s_insert;`, v.ReferenceTable, v.Name)

			qDropUniqueConstrainUpdate := fmt.Sprintf(`DROP INDEX IF EXISTS enforce_unique_%s_%s_update;`, v.ReferenceTable, v.Name)
			_, err := database_functions.ExecuteNonQuery(qDropIndex + qDropUniqueConstrainInsert + qDropUniqueConstrainUpdate)
			if err != nil {
				return err
			}
			break
		}
	}

	qRemoveColumn := fmt.Sprintf(`ALTER TABLE %s
		DROP COLUMN %s;`,
		tableName,
		columnConfig.Name)
	_, err = database_functions.ExecuteNonQuery(qRemoveColumn)
	if err != nil {
		return err
	}

	// update table config
	for i, v := range actualTableConfig.Columns {
		if v.Name == columnConfig.Name {
			actualTableConfig.Columns = append(actualTableConfig.Columns[:i], actualTableConfig.Columns[i+1:]...)
			break
		}
	}
	jsonStr, _ := tableConfigToJson(actualTableConfig)

	tablesConfig := database_config.GetTablesConfig()
	tablesConfig[tableName] = actualTableConfig
	database_config.SetTablesConfig(tablesConfig)

	qUpdateTableConfig := fmt.Sprintf(`UPDATE tables_config SET config = '%s' WHERE key = '%s';`, jsonStr, tableName)
	_, err = database_functions.ExecuteNonQuery(qUpdateTableConfig)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTable(tableName string) error {
	// check if table exists
	_, err := database_config.GetTableConfig(tableName)
	if err != nil {
		return err
	}

	qDeleteTable := fmt.Sprintf(`DROP TABLE %s;`, tableName)
	_, err = database_functions.ExecuteNonQuery(qDeleteTable)
	if err != nil {
		return err
	}

	qDeleteTableConfig := fmt.Sprintf(`DELETE FROM tables_config WHERE key='%s';`, tableName)
	_, err = database_functions.ExecuteNonQuery(qDeleteTableConfig)
	if err != nil {
		return err
	}

	qDeletePermissions := fmt.Sprintf(`DELETE FROM tables_permissions WHERE tableName = '%s';`, tableName)
	_, err = database_functions.ExecuteNonQuery(qDeletePermissions)
	if err != nil {
		return err
	}

	database_config.LoadTablesConfig()
	return nil
}

func tableConfigToJson(tableConf database_config.TableConfig) (string, error) {
	jsonString, err := json.Marshal(tableConf)
	if err != nil {
		return "", fmt.Errorf("failed to marshal/stringify column config: %v", err)
	}
	return string(jsonString), err
}
