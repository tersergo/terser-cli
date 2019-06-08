// Copyright (c) 2019 TerserGo
// 2019-05-20 10:42
// schema/table_schema.go

package schema

import "strings"

type TableSchema struct {
	DBName         string // database name
	Name           string // table name
	StructName     string // table model name
	VarName        string
	Comment        string         // table comment
	EngineName     string         // table engine
	FileName       string         // table file name
	ColumnList     []ColumnSchema // table columns
	PrimaryKeys    []ColumnSchema // primary columns list
	HasPrimaryKey  bool           // table primary key field
	IsIncrement    bool           // table is auto increment
	HasNullable    bool
	HasDateTime    bool
	HasEnum        bool
	LogicDeleteKey string
	CreateUserKey  string
	UpdateUserKey  string
	GenerateTime   string
	LabelTag       string
	LabelDot       string
}

func (t *TableSchema) SetIsIncrement(v interface{}) {
	i := toInt(v, 0)
	t.IsIncrement = i > 0
}

func (t *TableSchema) AppendColumn(column ColumnSchema) {
	if column.IsPrimaryKey {
		t.HasPrimaryKey = true
		column.Index = len(t.PrimaryKeys)
		t.PrimaryKeys = append(t.PrimaryKeys, column)
	}

	if !t.HasNullable && column.IsNullable {
		t.HasNullable = true
	}
	if !t.HasDateTime && column.IsDateTime {
		t.HasDateTime = true
	}
	if !t.HasNullable && column.IsNullable {
		t.HasNullable = true
	}
	if !t.HasEnum && column.IsEnum {
		t.HasEnum = true
	}

	switch strings.ToLower(column.Name) {
	case "is_deleted":
		t.LogicDeleteKey = column.PropertyName
	case "created_on":
		t.CreateUserKey = column.PropertyName
	case "updated_on":
		t.UpdateUserKey = column.PropertyName
	}

	t.ColumnList = append(t.ColumnList, column)
}

func (t *TableSchema) Init() {
	t.LabelTag = "`"
	t.LabelDot = "."
	t.LogicDeleteKey = ""
	t.CreateUserKey = ""
	t.UpdateUserKey = ""

	t.FileName = GetTableFileName(t.Name)
	t.StructName = GetFriendlyName(t.FileName)
	t.VarName = t.FileName[0:1]
}
