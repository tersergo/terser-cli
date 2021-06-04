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
	"github.com/tersergo/terser-cli/template/controller"
	"github.com/tersergo/terser-cli/template/model"
	"github.com/tersergo/terser-cli/template/proto"
)

var (
	dbName       string
	dsn          string
	appName      string
	dbDriver     string
	unsigned     string
	isInt32      string
	generateTime string
)

func init() {
	generateTime = time.Now().Format("2006/01/02 15:04:05")
	flag.StringVar(&dbName, "name", "test", "database name")
	flag.StringVar(&dsn, "dsn", "root:root@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local", "database source name")
	flag.StringVar(&appName, "app", "new-app", "app dir name")
	flag.StringVar(&dbDriver, "driver", "mysql", "database driver name: MySQL, SQLite3, Postgres, MSSQL")
	//flag.StringVar(&isInt32, "int32", "0", "default int=int32")
	flag.StringVar(&unsigned, "unsigned", "1", "support unsigned type(0=ignore): uint, uint8, uint16, uint32, uint64")
	flag.Parse()
}

func main() {
	prompt := `
terser-cli -help
terser-cli -name=test -dsn="root:root@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local" -driver=mysql -app=../new-app -unsigned=0
`
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

	for _, configTmpl := range modelFile {
		fmt.Printf("%#v", configTmpl)
		generateSingleFile(configTmpl.FileName, configTmpl.DirName, configTmpl.FileContent, query)
	}

	moduleModel := "model"
	tplModel, err := getTemplate(moduleModel, model.Model)
	if err != nil {
		fmt.Println("template parse error: ", moduleModel, err.Error())
		return
	}
	generateTemplateFile(tplModel, moduleModel, tableList)

	moduleProto := "proto"
	tplProto, err := getTemplate(moduleProto, proto.Proto)
	if err != nil {
		fmt.Println("template parse error: ", moduleProto, err.Error())
		return
	}
	generateTemplateFile(tplProto, moduleProto, tableList)

	moduleController := "controller"
	tplController, err := getTemplate(moduleController, controller.Controller)
	if err != nil {
		fmt.Println("template parse error: ", moduleController, err.Error())
		return
	}
	generateTemplateFile(tplController, moduleController, tableList)

}

var modelFile = []schema.FileTemplate{
	{FileName: "db_config.go", DirName: "model", FileContent: model.DBConfig},
	{FileName: "db_where.go", DirName: "model", FileContent: model.DBWhere},
	{FileName: "query_list.proto", DirName: "proto", FileContent: proto.QueryList},
}

func generateTemplateFile(modTpl *template.Template, moduleName string, tableList map[string]*schema.TableSchema) {
	for tableName, table := range tableList {
		table.GenerateTime = generateTime
		tpl, _ := modTpl.Clone()

		fileExt := ".go"
		if moduleName == "proto" {
			fileExt = ".proto"
		}

		appPath := fmt.Sprintf("%s/%s/%s%s", appName, moduleName, table.FileName, fileExt)
		fmt.Println(tableName, "\t\t=>\t", appPath)

		file, err := newFile(appName, moduleName, table.FileName+fileExt)
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
	tpl = template.New(name)
	tpl.Funcs(template.FuncMap{"Sum": Sum})
	tpl, err = tpl.Parse(content)
	if err == nil {
		tpl = template.Must(tpl, err)
	}

	return tpl, err
}

func Sum() func(nums ...int) (int, error) {
	return func(nums ...int) (int, error) {
		sum := 0
		for _, v := range nums {
			sum += v
		}
		return sum, nil
	}
}

func newFile(paths ...string) (*os.File, error) {
	pathDeep := len(paths) - 1
	createDirPath(paths[0:pathDeep]...)
	filePath := filepath.Join(paths...) //+ ".go"

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
