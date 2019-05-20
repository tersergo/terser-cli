package schema

import (
	"strings"
)

type ColumnSchema struct {
	TableName      string
	Name           string
	FriendlyName   string
	Comment        string
	ColumnType     string // 列具体类型 varchar(64), int(11)
	DataType       string // 数据类型(无精度) varchar, int
	DataTypeLength int    // 数据精度
	DataTypeScale  int    // 数据小数精度
	DefaultValue   string
	IsNullable     bool
	IsPrimaryKey   bool
	GoDataType     string // Golang对应的基础类型
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

func (c *ColumnSchema) InitGoDataType() {
	if len(c.GoDataType) > 0 {
		return
	}

}

func NewColumnSchema(list map[string]interface{}) (schema ColumnSchema) {
	if len(list) == 0 {
		return schema
	}

	/**
	TABLE_CATALOG	TABLE_SCHEMA	TABLE_NAME	COLUMN_NAME	ORDINAL_POSITION	COLUMN_DEFAULT	IS_NULLABLE	DATA_TYPE	CHARACTER_MAXIMUM_LENGTH	CHARACTER_OCTET_LENGTH	NUMERIC_PRECISION	NUMERIC_SCALE	DATETIME_PRECISION	CHARACTER_SET_NAME	COLLATION_NAME	COLUMN_TYPE	COLUMN_KEY	EXTRA	PRIVILEGES	COLUMN_COMMENT	GENERATION_EXPRESSION
	def	open_campus	t_sp_account	open_user_id	1		NO	varchar	64	256				utf8mb4	utf8mb4_general_ci	varchar(64)	PRI		select,insert,update,references	平台用户编号
	def	open_campus	t_sys_category	category_id	1		NO	int			10	0				int(10) unsigned	PRI	auto_increment	select,insert,update,references	分类编号
	*/
	schema = ColumnSchema{}

	for field, value := range list {
		switch strings.ToUpper(field) {
		case "TABLE_NAME":
			schema.TableName = toString(value)
		case "COLUMN_NAME":
			schema.Name = toString(value)
			schema.FriendlyName = GetFriendlyName(schema.Name)
		case "TABLE_COMMENT":
			schema.Comment = toString(value)
		case "COLUMN_DEFAULT":
			schema.DefaultValue = toString(value)
		case "COLUMN_KEY":
			schema.SetIsPrimaryKey(value)
		case "IS_NULLABLE":
			schema.SetIsNullable(value)
		case "COLUMN_TYPE":
			schema.ColumnType = toString(value)
		case "DATA_TYPE":
			schema.DataType = toString(value)
		case "CHARACTER_MAXIMUM_LENGTH", "NUMERIC_PRECISION":
			schema.SetDataTypeLength(value)
		case "NUMERIC_SCALE":
			schema.SetDataTypeScale(value)
		}
	}

	schema.InitGoDataType()

	return schema
}
