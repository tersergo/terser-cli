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

var (
	dbName       string
	dsn          string
	appName      string
	dbDriver     string
	unsigned     string
	generateTime string
)

func init() {
	generateTime = time.Now().Format("2006/01/02 15:04:05")
	flag.StringVar(&dbName, "name", "test", "database name")
	flag.StringVar(&dsn, "dsn", "root:root@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local", "database source name")
	flag.StringVar(&appName, "app", "new-app", "app dir name")
	flag.StringVar(&dbDriver, "driver", "mysql", "database driver name: MySQL, SQLite3, Postgres, MSSQL")
	flag.StringVar(&unsigned, "unsigned", "1", "support unsigned type(0=ignore): uint, uint8, uint16, uint32, uint64")
	flag.Parse()
}

func main() {
	prompt := `terser -name=test -dsn="root:root@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local" -driver=mysql -app=../new-app`
	if len(os.Args) == 1 {
		fmt.Println(prompt)
		os.Exit(1)
	}

	query := schema.NewQuery(dbName, dsn, strings.ToLower(dbDriver))
	schema.IgnoreUnsignedType = unsigned == "0"
	tableList, err := query.GetDBSchema()
	if err != nil {
		fmt.Println("database error: ", err.Error())
		return
	}
	dbConfigName := "db_config"
	generateSingleFile(dbConfigName, "model", tpl.DBConfig, query)

	moduleName := "model"
	tpl, err := getTemplate(moduleName, tpl.Model)
	if err != nil {
		fmt.Println("template parse error: ", moduleName, err.Error())
		return
	}
	generateTemplateFile(tpl, moduleName, tableList)

}

func generateTemplateFile(modTpl *template.Template, moduleName string, tableList map[string]*schema.TableSchema) {
	for tableName, table := range tableList {
		table.GenerateTime = generateTime
		tpl, _ := modTpl.Clone()

		appPath := fmt.Sprintf("%s/%s/%s.go", appName, moduleName, table.FileName)
		fmt.Println(tableName, "\t\t=>\t", appPath)

		file, err := newFile(appName, moduleName, table.FileName)
		if err == nil {
			err = tpl.Execute(file, table)
		}

		if err != nil {
			fmt.Println("template write error: ", err.Error())
		}

		file.Close()
	}
}

func getTemplate(name, content string) (tpl *template.Template, err error) {
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

func generateSingleFile(fileName, moduleName, tplConfig string, data interface{}) {
	tpl, err := getTemplate(fileName, tplConfig)
	if err != nil {
		fmt.Println("template parse error: ", fileName, err.Error())
		return
	}

	file, err := newFile(appName, moduleName, fileName)
	if err == nil {
		err = tpl.Execute(file, data)
	}
}

func createDirPath(paths ...string) {
	for index := 0; index < len(paths); index++ {
		length := index + 1
		fullPath := filepath.Join(paths[0:length]...)

		os.Mkdir(fullPath, os.ModePerm)
	}
}
