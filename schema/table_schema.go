package schema

import (
	"strings"
)

type TableSchema struct {
	DBName     string
	Name       string
	Comment    string
	EngineName string
	CreateTime string
	UpdateTime string

	FriendlyName string
	FileName     string
	IsIncrement  bool
	PrimaryKey   string

	ColumnList []ColumnSchema
}

func NewTableSchema(list map[string]interface{}) (schema TableSchema) {
	if len(list) == 0 {
		return schema
	}

	/**
	UPDATE `` SET `TABLE_CATALOG` = 'def', `TABLE_SCHEMA` = 'open_campus', `TABLE_NAME` = 't_sys_menu',
		`TABLE_TYPE` = 'BASE TABLE', `ENGINE` = 'InnoDB', `VERSION` = 10, `ROW_FORMAT` = 'Dynamic',
		`TABLE_ROWS` = 0, `AVG_ROW_LENGTH` = 0, `DATA_LENGTH` = 16384, `MAX_DATA_LENGTH` = 0, `INDEX_LENGTH` = 0,
		`DATA_FREE` = 0, `AUTO_INCREMENT` = 100, `CREATE_TIME` = '2019-05-19 23:28:52', `UPDATE_TIME` = NULL,
		`CHECK_TIME` = NULL, `TABLE_COLLATION` = 'utf8mb4_general_ci', `CHECKSUM` = NULL, `CREATE_OPTIONS` = '',
		`TABLE_COMMENT` = '系统菜单表'
	*/
	schema = TableSchema{}

	for field, value := range list {
		switch strings.ToUpper(field) {
		case "TABLE_SCHEMA":
			schema.DBName = toString(value)
		case "TABLE_NAME":
			schema.Name = toString(value)
			schema.FileName = GetTableFileName(schema.Name)
			schema.FriendlyName = GetFriendlyName(schema.FileName)
		case "TABLE_COMMENT":
			schema.Comment = toString(value)
		case "ENGINE":
			schema.EngineName = toString(value)
			schema.SetIsIncrement(value)
		case "CREATE_TIME":
			schema.CreateTime = toString(value)
		case "UPDATE_TIME":
			schema.UpdateTime = toString(value)
		}
	}

	return schema
}

func (t *TableSchema) SetIsIncrement(v interface{}) {
	i, err := toInt(v)
	if err == nil && i > 0 {
		t.IsIncrement = true
	} else {
		t.IsIncrement = false
	}
}
