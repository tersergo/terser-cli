package tpl

var DBConfigTemplate = `package model

import (
	"github.com/jinzhu/gorm"
)

// DBConfig 数据库配置类
type DBConfig struct {
	DBName     string
	DriverName string
	DSN        string // database source name
	dbServer   *gorm.DB
}


var defaultConfig *DBConfig

// 获取默认数据库配置
func DefaultConfig() (dc *DBConfig) {
	if defaultConfig == nil {
		defaultConfig = &DBConfig{
			DBName:     "{{.DBName}}",
			DriverName: "{{if .DriverName}}{{.DriverName}}{{else}}mysql{{end}}",
			DSN:        "{{.DSN}}",
		}
	}

	return defaultConfig
}

func (dc *DBConfig) initDB() {
	if dc.dbServer != nil {
		return
	}

	db, err := gorm.Open(dc.DriverName, dc.DSN)
	if err != nil {
		// todo: error handle
		panic(err)
	}
	dc.dbServer = db
}

func (dc *DBConfig) GetDBServer() (db *gorm.DB) {
	dc.initDB()
	return dc.dbServer
}

func (dc *DBConfig) Close() (err error) {
	if dc.dbServer != nil {
		err = dc.dbServer.Close()
	}

	return err
}

`
