// Copyright (c) 2019 TerserGo
// 2019-05-20 10:42
// schema/table_schema.go

package schema

type TableSchema struct {
	DBName         string // database name
	Name           string // table name
	Comment        string // table comment
	EngineName     string // table engine
	ModelName      string // table model name
	FileName       string // table file name
	IsIncrement    bool   // table is auto increment
	HasPrimaryKey  bool   // table primary key field
	HasNullable    bool
	LogicDeleteKey string
	ColumnList     []ColumnSchema // table columns
	primaryKeyIds  []int
	CreateTime     string
	ShortName      string
}

func (t *TableSchema) SetIsIncrement(v interface{}) {
	i, err := toInt(v)
	if err == nil && i > 0 {
		t.IsIncrement = true
	} else {
		t.IsIncrement = false
	}
}

func (t *TableSchema) AppendColumn(column ColumnSchema) {
	if column.IsPrimaryKey {
		t.HasPrimaryKey = true
		t.primaryKeyIds = append(t.primaryKeyIds, len(t.ColumnList))
	}

	if !t.HasNullable && column.IsNullable {
		t.HasNullable = true
	}

	t.ColumnList = append(t.ColumnList, column)
}

func (t *TableSchema) GetPrimaryKeys() (columns []ColumnSchema) {
	if len(t.primaryKeyIds) > 0 {
		for _, index := range t.primaryKeyIds {
			columns = append(columns, t.ColumnList[index])
		}
	}
	return columns
}

func (t *TableSchema) InitName() {
	t.FileName = GetTableFileName(t.Name)
	t.ModelName = GetFriendlyName(t.FileName)
	t.ShortName = t.FileName[0:1]
}
