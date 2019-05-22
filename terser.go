package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/tersergo/terser-cli/schema"
	tpl "github.com/tersergo/terser-cli/template"
)

func main() {

	help := `terser -name=dbName -dsn="root:root@tcp(localhost:3306)/test?parseTime=true" -driver=mysql -app=new-app`

	if len(os.Args) == 1 {
		fmt.Println(help)
		os.Exit(0)
	}

	var (
		dbName   string
		dsn      string
		appName  string
		dbDriver string
	)

	flag.StringVar(&dbName, "name", "test", "database name")
	flag.StringVar(&dsn, "dsn", "root:root@tcp(localhost:3306)/test?parseTime=true", "database source name")
	flag.StringVar(&appName, "app", "new-app", "app dir name")
	flag.StringVar(&dbDriver, "driver", "mysql", "database driver name: MySQL, SQLite, MariaDB")

	flag.Parse()
	modName := "model"

	query := schema.NewQuery(dbName, dsn, strings.ToLower(dbDriver))

	tableList, err := query.GetDBSchema()
	if err != nil {
		fmt.Println("database error: ", err.Error())
		return
	}

	dbConfigName := "db_config"
	configTpl, err := generateTemplate(dbConfigName, tpl.DBConfigTemplate)
	if err != nil {
		fmt.Println("template parse error: ", dbConfigName, err.Error())
		return
	}

	file, err := newFile(appName, "model", dbConfigName)
	if err == nil {
		err = configTpl.Execute(file, query)
	}

	tplName := modName
	modelTpl, err := generateTemplate(tplName, tpl.ModelTemplate)
	if err != nil {
		fmt.Println("template parse error: ", tplName, err.Error())
		return
	}
	generateTime := time.Now().Format("2006/01/02 15:04:05")

	for tableName, table := range tableList {
		table.GenerateTime = generateTime
		tpl, _ := modelTpl.Clone()
		fmt.Println(tableName, "\t\t=>\t", table.StructName)

		//var stream bytes.Buffer
		file, err := newFile(appName, modName, table.FileName)
		if err == nil {
			err = tpl.Execute(file, table)
		}
		if err != nil {
			fmt.Println("template write error: ", err.Error())
		}

		file.Close()
	}
}

func createDirPath(paths ...string) {
	for index := 0; index < len(paths); index++ {
		length := index + 1
		fullPath := filepath.Join(paths[0:length]...)

		os.Mkdir(fullPath, os.ModePerm)
	}
}

func generateTemplate(name, content string) (tpl *template.Template, err error) {
	tpl, err = template.New(name).Parse(content)
	if err == nil {
		tpl = template.Must(tpl, err)
	}
	return tpl, err
}

func newFile(paths ...string) (*os.File, error) {
	pathDeep := len(paths) - 1
	createDirPath(paths[0:pathDeep]...)
	filePath := filepath.Join(paths...) + ".go"

	return os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
}

//func writeFile(dir, file string, stream bytes.Buffer) {
//	os.Mkdir("../"+dir, os.ModePerm)
//
//	filePath := filepath.Join("../", dir, file+".go")
//	data, err := format.Source(stream.Bytes())
//
//	if err != nil {
//		fmt.Println("stream error: " + err.Error())
//		return
//	}
//	ioutil.WriteFile(filePath, data, os.ModePerm)
//}

//func generateDBConfigFile(query *schema.Query) {
//
//	dbTemp, err := template.New("db_config").Parse(tpl.DBConfigTemplate)
//	if err != nil {
//		fmt.Println("template error: ", err.Error())
//	}
//	tpl := template.Must(dbTemp, err)
//
//	var stream bytes.Buffer
//	err = tpl.Execute(&stream, query)
//	if err != nil {
//		fmt.Println("db_config template error: ", err.Error())
//	}
//
//	writeFile("model", "db_config", stream)
//
//}
