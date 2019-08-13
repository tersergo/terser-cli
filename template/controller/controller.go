package controller

var Controller = `package controller
type Create{{.StructName}} struct { 
{{range .ColumnList}}	{{.PropertyName}}	{{.GoDataType}}		{{$.LabelTag}}form:"{{.PropertyName}}" json:"{{.PropertyName}}" required:"false"{{$.LabelTag}}	// {{.Comment}}{{if .IsNullable}} [null type]{{end}}{{if .DefaultValue}} [default: {{.DefaultValue}}]{{end}}
{{end}}}

type Update{{.StructName}} struct { 
{{range .ColumnList}}	{{.PropertyName}}	{{.GoDataType}}		{{$.LabelTag}}form:"{{.PropertyName}}" json:"{{.PropertyName}}" required:"false"{{$.LabelTag}}	// {{.Comment}}{{if .IsNullable}} [null type]{{end}}{{if .DefaultValue}} [default: {{.DefaultValue}}]{{end}}
{{end}}}

`
