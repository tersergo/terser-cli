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
	splitChars = [...]string{"_", "-"}

	tablePrefixes = [...]string{"t_"}

	friendlyNameMaps = map[string]string{
		"db":   "DB",
		"id":   "ID",
		"pk":   "PK",
		"uuid": "UUID",
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
		"tinyint":    "int",
		"smallint":   "int",
		"mediumint":  "int",
		"int":        "int32",
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
	if len(vs) == 0 {
		return false
	}

	return vs == s
}

func toInt(v interface{}) (num int, err error) {
	i, ok := v.(int)
	if ok {
		return i, nil
	}

	s := toString(v)
	if len(s) > 0 {
		i, err := strconv.Atoi(s)
		if err == nil && i > 0 {
			num = i
		}
	} else {
		err = fmt.Errorf("values is null")
	}
	return num, err
}

func GetFriendlyName(original string) (friendlyName string) {
	if len(original) == 0 {
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
		length := len(name)

		if length == 0 {
			continue
		}

		if length == 1 {
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
