// Copyright (c) 2019 TerserGo
// 2019-05-20 10:42
// schema/column_schema.go

package schema

import (
	"strings"
)

type ColumnSchema struct {
	TableName      string
	Name           string // 表的字段名
	PropertyName   string // struct属性名
	VarName        string // func变量名
	Comment        string
	ColumnType     string // 列具体数据类型 varchar(64), int(11)
	DataType       string // 数据类型(无精度) varchar, int
	DataTypeLength int    // 数据精度
	DataTypeScale  int    // 数据小数精度
	DefaultValue   string
	IsNumeral      bool
	IsDateTime     bool
	IsNullable     bool
	IsPrimaryKey   bool
	IsEnum         bool
	IsJson         bool
	GoDataType     string // Golang对应的基础类型
	Index          int
}

func (c *ColumnSchema) SetIsPrimaryKey(v interface{}) {
	c.IsPrimaryKey = equalToString(v, "PRI")
}

func (c *ColumnSchema) SetIsNullable(v interface{}) {
	c.IsNullable = equalToString(v, "YES")
}

func (c *ColumnSchema) SetDataTypeLength(v interface{}) {
	i, err := toInt(v)
	if err == nil && i > 0 {
		c.DataTypeLength = i
	}
}

func (c *ColumnSchema) SetDataTypeScale(v interface{}) {
	i, err := toInt(v)
	if err == nil && i > 0 {
		c.DataTypeScale = i
	}
}

func (c *ColumnSchema) Init() {
	if len(c.GoDataType) > 0 {
		return
	}

	if len(c.PropertyName) > 0 {
		_, exist := friendlyNameMaps[c.Name]
		if !exist {
			c.VarName = strings.ToLower(c.PropertyName[0:1]) + c.PropertyName[1:]
		} else {
			c.VarName = strings.ToLower(c.Name)
		}
	}

	goBaseType, exists := dataTypeMaps[c.DataType]

	if !exists {
		c.GoDataType = c.DataType
		return
	}
	c.GoDataType = goBaseType

	switch strings.ToLower(c.DataType) {
	case "tinyint", "smallint", "mediumint", "int", "bigint":
		c.IsNumeral = true
		if c.IsNullable {
			c.GoDataType = "null.Int"
		} else if strings.Index(c.ColumnType, "unsigned") > 0 {
			c.ColumnType = strings.Replace(c.ColumnType, " unsigned", "", 1)
			if !IgnoreUnsignedType {
				c.GoDataType = "u" + goBaseType
			}
		}
	case "date", "year", "time", "timestamp", "datetime":
		c.IsDateTime = true
		if strings.Index(c.DefaultValue, "CURRENT_TIMESTAMP") >= 0 {
			c.DefaultValue = ""
		}
		if c.IsNullable {
			c.GoDataType = "null.Time"
		}
	case "float", "double", "decimal":
		c.IsNumeral = true
		if c.IsNullable {
			c.GoDataType = "null.Float"
		}
	case "bool":
		c.IsNumeral = true
		if c.IsNullable {
			c.GoDataType = "null.Bool"
		}
	case "enum":
		c.IsEnum = true
		//c.initEnumType(c.ColumnType)
	case "char", "varchar", "text", "json", "tinytext", "mediumtext", "longtext":
		c.IsJson = strings.ToLower(c.DataType) == "json"
		if c.IsNullable {
			c.GoDataType = "null.String"
		}
	}

}

//func (c *ColumnSchema) initEnumType(enumType string) {
//
//}
