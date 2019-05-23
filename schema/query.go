// Copyright (c) 2019 TerserGo
// 2019-05-20 10:42
// schema/query.go

package schema

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Query struct {
	DBName     string
	DSN        string
	DriverName string
	db         *sql.DB
}

func NewQuery(dbName, dsn, dbDriver string) (dbQuery *Query) {
	return &Query{
		DBName:     dbName,
		DSN:        dsn,
		DriverName: dbDriver,
	}
}

func (query *Query) initDB() {
	if query.db != nil {
		return
	}
	db, err := sql.Open(query.DriverName, query.DSN)

	if err != nil {
		fmt.Println("db error: ", err.Error())
		panic(err)
	}

	query.db = db
}

func (query *Query) GetDBSchema() (tableMaps map[string]*TableSchema, err error) {
	query.initDB()
	defer query.Close()

	tableList, err := query.fetchTableSchema()
	tableMaps = make(map[string]*TableSchema, len(tableList))

	for _, table := range tableList {
		tableMaps[table.Name] = table
	}

	columnList, err := query.fetchColumnSchema()

	for _, column := range columnList {
		name := column.TableName
		table, exists := tableMaps[name]
		if exists {
			table.AppendColumn(column)
		}
	}

	return tableMaps, err
}

func (query *Query) fetchTableSchema() (tableList []*TableSchema, err error) {
	sql := "SELECT TABLE_SCHEMA, TABLE_NAME, TABLE_COMMENT, `ENGINE`, AUTO_INCREMENT FROM information_schema.`TABLES` WHERE TABLE_SCHEMA = ?"
	rows, err := query.db.Query(sql, query.DBName)
	if err != nil {
		return tableList, err
	}

	for rows.Next() {
		var (
			dbName                = ""
			tbName                = ""
			tbComment             = ""
			tbEngine              = ""
			tbInc     interface{} = 0
		)
		subErr := rows.Scan(&dbName, &tbName, &tbComment, &tbEngine, &tbInc)
		if subErr != nil {
			fmt.Println("table query error: ", subErr.Error())
		}

		table := &TableSchema{
			DBName:     dbName,
			Name:       tbName,
			Comment:    tbComment,
			EngineName: tbEngine,
		}

		table.SetIsIncrement(tbInc)
		table.Init()

		tableList = append(tableList, table)
	}

	return tableList, err
}

func (query *Query) fetchColumnSchema() (columnList []ColumnSchema, err error) {
	sql := `SELECT 
       TABLE_NAME, COLUMN_NAME, COLUMN_COMMENT, COLUMN_DEFAULT, COLUMN_KEY, IS_NULLABLE, COLUMN_TYPE, 
       DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE 
		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ?  ORDER BY ORDINAL_POSITION `
	rows, err := query.db.Query(sql, query.DBName)
	if err != nil {
		return columnList, err
	}
	for rows.Next() {
		var (
			tbName                  = ""
			colName                 = ""
			colComment              = ""
			colDefault  interface{} = ""
			colKey                  = ""
			colNullable             = ""
			colShowType             = ""
			colDataType             = ""
			charLength  interface{} = 0
			numLength   interface{} = 0
			scaleLength interface{} = 0
		)
		subErr := rows.Scan(&tbName, &colName, &colComment, &colDefault, &colKey, &colNullable, &colShowType,
			&colDataType, &charLength, &numLength, &scaleLength)

		if subErr != nil {
			fmt.Println("column query error: ", subErr.Error())
		}

		column := ColumnSchema{
			TableName:    tbName,
			Name:         colName,
			Comment:      colComment,
			DefaultValue: toString(colDefault),
			ColumnType:   colShowType,
			DataType:     colDataType,
		}

		column.PropertyName = GetFriendlyName(colName)
		column.SetIsPrimaryKey(colKey)
		column.SetIsNullable(colNullable)
		column.SetDataTypeLength(charLength)
		column.SetDataTypeLength(numLength)
		column.SetDataTypeScale(scaleLength)

		column.Init()

		columnList = append(columnList, column)
	}

	return columnList, err
}

func (query *Query) Close() {
	if query.db != nil {
		query.db.Close()
	}

	query.db = nil
}
