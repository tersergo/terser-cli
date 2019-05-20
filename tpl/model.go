package tpl

var ModelTemplate = `// auto-generated by terser-cli
// table: {{.Name}} {{.CreateTime}}
package model


import (
	"github.com/jinzhu/gorm"
	"time"
//  "database/sql"
//  "github.com/guregu/null"
)

// {{.Comment}}
type {{.ModelName}} struct {
	//gorm.Model // 使用gorm基础model(包含ID,CreatedAt,UpdatedAt,DeletedAt字段定义)
	{{range $index, $column .= ColumnList}}{{$column.VarName}}	{{$column.GoDataType}}		{{$.LabelTag}}gorm:"column:{{.Name}};type:{{.ColumnType}};{{if .IsPrimaryKey}}primary_key{{end}}" json:"{{.Name}}"{{$.LabelTag}}	// {{.Comment}} {{if .DefaultValue}}默认{{.DefaultValue}}{{end}}
{{end}}
}

// {{.ModelName}}'s table name
func ({{.ShortName}} *{{.ModelName}}) TableName() string {
	return "{{.Name}}"
}

//func AutoMigrate(db *gorm.DB) {
//	db.AutoMigrate(&{{.ModelName}}{})
//}

//func ({{.ShortName}} *{{.ModelName}}) BeforeCreate(scope *gorm.Scope) error {
//  return nil
//}

`
