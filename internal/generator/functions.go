package generator

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/isacikgoz/mattermost-suite-utilities/internal/model"
)

var funcMap = template.FuncMap{
	// make s string start with upper case
	"public": func(s string) string {
		return strings.Title(s)
	},
	"sliceTypes": func(s []*model.Field) map[string]string {
		ret := map[string]string{}
		for _, r := range s {
			if strings.Contains(r.Type, "[]") {
				ret["SliceOf"+strings.Title(r.Type[2:])] = r.Type[2:]
			}
		}
		return ret
	},
	"customType": func(s string) bool {
		switch s {
		case "bool", "byte", "[]byte", "string", "[]string", "int", "int32", "int64", "float32", "float64":
			return false
		default:
			return true
		}
	},
	"generateInitializer": func(t, name, ft string) string {
		s := fmt.Sprintf("%s.%s", strings.ToLower(string(t[0])), strings.Title(name))
		if strings.Contains(ft, "[]") {
			return fmt.Sprintf("NewSliceOf%s(%s)", strings.Title(ft[2:]), s)
		}
		return s
	},
	"generateSetStatement": func(t string) string {
		if strings.Contains(t, "[]") {
			return ".Replace(v)"
		}
		return " = v"
	},
	"generateGetStatement": func(t string) string {
		if strings.Contains(t, "[]") {
			return ".Range()"
		}
		return ""
	},
	"processType": func(s string) string {
		if strings.Contains(s, "[]") {
			return "SliceOf" + strings.Title(s[2:])
		}
		return s
	},
}

func patch_field(tags map[string][]string) bool {
	for k, v := range tags {
		if k == "model" {
			for _, s := range v {
				if s == "patch" {
					return true
				}
			}
		}
	}
	return false
}
