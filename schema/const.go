// Copyright (c) 2019 TerserGo
// 2019-05-20 10:42
// schema/const.go

package schema

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	IgnoreUnsignedType = false
	IsGRPCModel        = false

	splitChars = [...]string{"_", "-"}

	tablePrefixes = [...]string{"t_", "a_"}

	friendlyNameMaps = map[string]string{
		"db":   "DB",
		"id":   "Id",
		"pk":   "PK",
		"uuid": "UUID",
		"ip":   "IP",
	}

	dataTypeMaps = map[string]string{
		"char":       "string",
		"varchar":    "string",
		"json":       "string",
		"text":       "string",
		"tinytext":   "string",
		"mediumtext": "string",
		"longtext":   "string",
		"date":       "time.Time",
		"year":       "time.Time",
		"time":       "time.Time",
		"timestamp":  "time.Time",
		"datetime":   "time.Time",
		"tinyint":    "int32",
		"smallint":   "int32",
		"mediumint":  "int32",
		"int":        "int32",
		//"tinyint":    "int8",
		//"smallint":   "int16",
		//"mediumint":  "int16",
		//"int":        "int",
		"bigint":     "int64",
		"float":      "float32",
		"double":     "float64",
		"decimal":    "float64",
		"enum":       "string",
		"bool":       "bool",
		"set":        "string",
		"blob":       "[]byte",
		"tinyblob":   "[]byte",
		"mediumblob": "[]byte",
		"longblob":   "[]byte",
		"binary":     "[]byte",
		"varbinary":  "[]byte",
	}

	nullTypeMaps = map[string]string{}
)

func toString(v interface{}) string {
	if v == nil {
		return ""
	}

	text, ok := v.(string)
	if ok {
		return text
	}

	return fmt.Sprintf("%s", v)
}

func equalToString(v interface{}, s string) bool {
	vs := toString(v)

	return vs == s
}

func toInt(v interface{}, defaultReturn int) (num int) {
	i, ok := v.(int)
	if ok {
		return i
	}

	s := toString(v)
	if s != "" {
		i, err := strconv.Atoi(s)
		if err == nil {
			return i
		}
	}

	return defaultReturn
}

func GetFriendlyName(original string) (friendlyName string) {
	if original == "" {
		return original
	}
	name := strings.ToLower(original)

	for _, split := range splitChars {
		if strings.Index(name, split) >= 0 {
			splitNames := strings.Split(original, split)
			return GetHumpName(splitNames...)
		}
	}

	return GetHumpName(original)
}

func GetHumpName(names ...string) (humpName string) {
	for _, name := range names {
		if name == "" {
			continue
		}

		if len(name) == 1 {
			humpName += strings.ToUpper(name)
			continue
		}

		fixName, hasName := friendlyNameMaps[name]
		if hasName {
			humpName += fixName
			continue
		}

		humpName += strings.ToUpper(name[0:1]) + name[1:]
	}

	return humpName
}

func GetTableFileName(original string) (name string) {
	name = strings.ToLower(original)

	for _, prefix := range tablePrefixes {
		if strings.Index(name, prefix) == 0 {
			startIndex := len(prefix)
			return name[startIndex:]
		}
	}

	return name

}
