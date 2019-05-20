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
		"db": "DB",
		"id": "ID",
	}

	dataTypeMaps = map[string]string{
		"char":       "string",
		"varchar":    "string",
		"json":       "string",
		"text":       "string",
		"mediumtext": "string",
		"longtext":   "string",
		"time":       "time.time",
		"timestamp":  "time.time",
		"datetime":   "time.time",
		"tinyint":    "int",
		"smallint":   "int",
		"int":        "int",
		"bigint":     "int",
		"float":      "float",
		"double":     "float",
		"decimal":    "float",
		"set":        "",
		"enum":       "",
		"blob":       "",
		"mediumblob": "",
		"longblob":   "",
	}
)

func toString(v interface{}) string {
	s, ok := v.(string)
	if ok {
		return s
	}

	return fmt.Sprint(v)
}

func equalToString(v interface{}, s string) bool {
	vs := toString(v)
	if len(vs) == 0 {
		return false
	}

	return vs == s
}

func toInt(v interface{}) (num int, err error) {
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
			return GetHumpName(splitNames)
		}
	}

	return GetHumpName([]string{original})
}

func GetHumpName(names []string) (humpName string) {
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
			startIndex := len(prefix) - 1
			return name[startIndex:]
		}
	}

	return name

}
