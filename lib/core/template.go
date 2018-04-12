/*
 * Copyright 2018 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package core

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
)

var (
	FuncMap = map[string]interface{}{
		"hello":      func() string { return "hello!" },
		"env":        func(s string) string { return os.Getenv(s) },
		"expandenv":  func(s string) string { return os.ExpandEnv(s) },
		"base":       path.Base,
		"dir":        path.Dir,
		"clean":      path.Clean,
		"ext":        path.Ext,
		"isAbs":      path.IsAbs,
		"quote":      quote,
		"squote":     squote,
		"contains":   func(substr string, str string) bool { return strings.Contains(str, substr) },
		"hasPrefix":  func(substr string, str string) bool { return strings.HasPrefix(str, substr) },
		"hasSuffix":  func(substr string, str string) bool { return strings.HasSuffix(str, substr) },
		"trim":       strings.TrimSpace,
		"trimAll":    func(a, b string) string { return strings.Trim(b, a) },
		"trimSuffix": func(a, b string) string { return strings.TrimSuffix(b, a) },
		"trimPrefix": func(a, b string) string { return strings.TrimPrefix(b, a) },
		"split":      split,
		"splitList":  func(sep, orig string) []string { return strings.Split(orig, sep) },
		"toString":   strval,
		"toStrings":  strslice,
		"join":       join,
		"sortAlpha":  sortAlpha,
		"b64enc":     base64encode,
		"b64dec":     base64decode,
		"cat":        cat,
		"indent":     indent,
		"nindent":    nindent,
		"replace":    replace,
		"atoi":       func(a string) int { i, _ := strconv.Atoi(a); return i },
		"int64":      toInt64,
		"int":        toInt,
		"float64":    toFloat64,
		"toJson":     toJson,
		"toYaml":     toYaml,
		"max":        max,
		"min":        min,
		"typeOf":     typeOf, // Reflection
		"typeIs":     typeIs,
		"typeIsLike": typeIsLike,
		"kindOf":     kindOf,
		"kindIs":     kindIs,
		"list":       list,   // Data Structures:
		"timeOf":     timeOf, // time to second

	}
)

var (
	parser *template.Template
)

func init() {
	parser = template.New("parser").Funcs(FuncMap)
}

func ParseTemplateFile(values interface{}, filename string) ([]byte, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	return ParseTemplateText(values, string(b))
}

func ParseTemplateText(values interface{}, text string) ([]byte, error) {
	var b bytes.Buffer

	tpl, err := parser.Parse(text)
	if err != nil {
		return b.Bytes(), err
	}
	err = tpl.Execute(&b, values)
	return b.Bytes(), err
}

func quote(str ...interface{}) string {
	out := make([]string, len(str))
	for i, s := range str {
		out[i] = fmt.Sprintf("%q", strval(s))
	}
	return strings.Join(out, " ")
}

func squote(str ...interface{}) string {
	out := make([]string, len(str))
	for i, s := range str {
		out[i] = fmt.Sprintf("'%v'", s)
	}
	return strings.Join(out, " ")
}

func cat(v ...interface{}) string {
	r := strings.TrimSpace(strings.Repeat("%v ", len(v)))
	return fmt.Sprintf(r, v...)
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

func nindent(spaces int, v string) string {
	return "\n" + indent(spaces, v)
}

func replace(old, new, src string) string {
	return strings.Replace(src, old, new, -1)
}

func join(sep string, v interface{}) string {
	return strings.Join(strslice(v), sep)
}

func split(sep, orig string) map[string]string {
	parts := strings.Split(orig, sep)
	res := make(map[string]string, len(parts))
	for i, v := range parts {
		res["_"+strconv.Itoa(i)] = v
	}
	return res
}

func sortAlpha(list interface{}) []string {
	k := reflect.Indirect(reflect.ValueOf(list)).Kind()
	switch k {
	case reflect.Slice, reflect.Array:
		a := strslice(list)
		s := sort.StringSlice(a)
		s.Sort()
		return s
	}
	return []string{strval(list)}
}

func strslice(v interface{}) []string {
	switch v := v.(type) {
	case []string:
		return v
	case []interface{}:
		l := len(v)
		b := make([]string, l)
		for i := 0; i < l; i++ {
			b[i] = strval(v[i])
		}
		return b
	default:
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Array, reflect.Slice:
			l := val.Len()
			b := make([]string, l)
			for i := 0; i < l; i++ {
				b[i] = strval(val.Index(i).Interface())
			}
			return b
		default:
			return []string{strval(v)}
		}
	}
}

func strval(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// toFloat64 converts 64-bit floats
func toFloat64(v interface{}) float64 {
	if str, ok := v.(string); ok {
		iv, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0
		}
		return iv
	}

	val := reflect.Indirect(reflect.ValueOf(v))
	switch val.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return float64(val.Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return float64(val.Uint())
	case reflect.Uint, reflect.Uint64:
		return float64(val.Uint())
	case reflect.Float32, reflect.Float64:
		return val.Float()
	case reflect.Bool:
		if val.Bool() == true {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func toInt(v interface{}) int {
	//It's not optimal. Bud I don't want duplicate toInt64 code.
	return int(toInt64(v))
}

// toInt64 converts integer types to 64-bit integers
func toInt64(v interface{}) int64 {
	if str, ok := v.(string); ok {
		iv, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return 0
		}
		return iv
	}

	val := reflect.Indirect(reflect.ValueOf(v))
	switch val.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return val.Int()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return int64(val.Uint())
	case reflect.Uint, reflect.Uint64:
		tv := val.Uint()
		if tv <= math.MaxInt64 {
			return int64(tv)
		}
		// TODO: What is the sensible thing to do here?
		return math.MaxInt64
	case reflect.Float32, reflect.Float64:
		return int64(val.Float())
	case reflect.Bool:
		if val.Bool() == true {
			return 1
		}
		return 0
	default:
		return 0
	}
}
func max(a interface{}, i ...interface{}) int64 {
	aa := toInt64(a)
	for _, b := range i {
		bb := toInt64(b)
		if bb > aa {
			aa = bb
		}
	}
	return aa
}

func min(a interface{}, i ...interface{}) int64 {
	aa := toInt64(a)
	for _, b := range i {
		bb := toInt64(b)
		if bb < aa {
			aa = bb
		}
	}
	return aa
}

func toJson(v interface{}) string {
	output, _ := json.Marshal(v)
	return string(output)
}

func toYaml(v interface{}) string {
	output, _ := yaml.Marshal(v)
	return string(output)
}

func base64encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func base64decode(v string) string {
	data, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// typeIs returns true if the src is the type named in target.
func typeIs(target string, src interface{}) bool {
	return target == typeOf(src)
}

func typeIsLike(target string, src interface{}) bool {
	t := typeOf(src)
	return target == t || "*"+target == t
}

func typeOf(src interface{}) string {
	return fmt.Sprintf("%T", src)
}

func kindIs(target string, src interface{}) bool {
	return target == kindOf(src)
}

func kindOf(src interface{}) string {
	return reflect.ValueOf(src).Kind().String()
}

func list(v ...interface{}) []interface{} {
	return v
}

func timeOf(v string) string {
	var n int64

	switch byte(v[len(v)-1]) {
	case 's', 'S':
		n = toInt64(v[:len(v)-1])
	case 'm', 'M':
		n = toInt64(v[:len(v)-1]) * 60
	case 'h', 'H':
		n = toInt64(v[:len(v)-1]) * 3600
	case 'd', 'D':
		n = toInt64(v[:len(v)-1]) * 3600 * 24
	default:
		n = toInt64(v)
	}
	return fmt.Sprintf("%d", n)
}
