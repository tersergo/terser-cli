package schema

import (
	"database/sql"
	"fmt"
)

type DBType string

type DBQuery struct {
	DriverName string
	DBName     string
	DSN        string

	db *sql.DB
}

func NewDBQuery(dbName, dsn string) (dbQuery *DBQuery) {
	return &DBQuery{
		DriverName: "mysql",
		DBName:     dbName,
		DSN:        dsn,
	}
}

func (d *DBQuery) initDb() {
	if d.db != nil {
		return
	}
	db, err := sql.Open(d.DriverName, d.DSN)

	if err != nil {
		fmt.Println("db error: ", err.Error())
		panic(err)
	}

	d.db = db
}

func (d *DBQuery) QueryList(sql string, arguments ...interface{}) (list []map[string]interface{}, err error) {
	d.initDb()

	rows, err := d.db.Query(sql, arguments)

	if err != nil {
		return list, err
	}

	for rows.Next() {
		maps := make(map[string]interface{})

		err = rows.Scan(&maps)

		list = append(list, maps)
	}

	return list, err
}

func (d *DBQuery) GetTableSchema() (tableMaps map[string]TableSchema, err error) {
	tableSql := "SELECT * FROM information_schema.`TABLES` WHERE TABLE_SCHEMA = ?"

	schemaList, err := d.QueryList(tableSql, d.DBName)
	if err != nil {
		return tableMaps, err
	}

	tableMaps = make(map[string]TableSchema, len(schemaList))

	for _, schema := range schemaList {
		table := NewTableSchema(schema)

		if len(table.Name) > 0 {
			tableMaps[table.Name] = table
		}
	}

	columnList, err := d.GetColumnSchema()

	for _, column := range columnList {

		tableName := column.TableName

	}

	return tableMaps, err
}

func (d *DBQuery) GetColumnSchema() (columnList []ColumnSchema, err error) {
	columnSql := " SELECT * FROM information_schema.`COLUMNS` WHERE TABLE_SCHEMA = ?"

	schemaList, err := d.QueryList(columnSql, d.DBName)
	if err != nil {
		return columnList, err
	}

	columnList = make([]ColumnSchema, len(schemaList))

	for index, schema := range schemaList {
		column := NewTableSchema(schema)

		if len(column.Name) > 0 {
			columnList[index] = NewColumnSchema(schema)
		}
	}

	return columnList, err
}

func (d *DBQuery) Close() {
	if d.db != nil {
		d.db.Close()
	}

	d.db = nil
}
