// Copyright (c) 2019 TerserGo
// 2019-05-20 10:42
// schema/table_schema.go

package schema

import "strings"

type TableSchema struct {
	DBName          string // database name
	Name            string // table name
	StructName      string // table model name
	VarName         string
	Comment         string // table comment
	EngineName      string // table engine
	FileName        string // table file name
	IsIncrement     bool   // table is auto increment
	HasNullable     bool
	LogicDeleteKey  string
	CreateUserKey   string
	UpdateUserKey   string
	ColumnList      []ColumnSchema // table columns
	PrimaryKeyCount int            // table primary key field
	PrimaryKeys     []ColumnSchema
	GenerateTime    string
	LabelTag        string
	LabelDot        string
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
		t.PrimaryKeys = append(t.PrimaryKeys, column)
		t.PrimaryKeyCount = len(t.PrimaryKeys)
	}

	if !t.HasNullable && column.IsNullable {
		t.HasNullable = true
	}

	switch strings.ToLower(column.Name) {
	case "is_deleted":
		t.LogicDeleteKey = column.PropertyName
	case "created_on":
		t.CreateUserKey = column.PropertyName
	case "modified_on":
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
