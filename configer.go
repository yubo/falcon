// Copyright 2017 Xiaomi, Inc.
// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package falcon

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	APP_CONF_DEFAULT = iota
	APP_CONF_DB
	APP_CONF_FILE
	APP_CONF_SIZE
)

var (
	APP_CONF_NAME = [APP_CONF_SIZE]string{
		"default", "db", "file",
	}
)

type Configer struct {
	data [APP_CONF_SIZE]map[string]string
}

func (c Configer) String() string {
	s := ""
	for i := 0; i < len(c.data); i++ {
		s1 := ""
		for k, v := range c.data[i] {
			s1 += fmt.Sprintf("%-17s %s\n", k, v)
		}
		if len(s1) > 0 {
			s += fmt.Sprintf("%-17s {\n%s\n}\n",
				APP_CONF_NAME[i], IndentLines(1, s1))
		} else {
			s += fmt.Sprintf("%-17s { }\n",
				APP_CONF_NAME[i])
		}
	}
	return s
}

func (c *Configer) Set(model int, m map[string]string) error {
	if model >= APP_CONF_SIZE || model < 0 {
		return errors.New("no model")
	}
	data := make(map[string]string)
	for k, v := range m {
		if len(k) == 0 {
			return errors.New("empty key")
		}
		data[strings.ToLower(k)] = v
	}
	c.data[model] = data
	return nil
}

func (c Configer) Get() [APP_CONF_SIZE]map[string]string {
	return c.data
}

// Bool returns the boolean value for a given key.
func (c *Configer) Bool(key string) (bool, error) {
	return ParseBool(c.getdata(key))
}

// DefaultBool returns the boolean value for a given key.
// if err != nil return defaltval
func (c *Configer) DefaultBool(key string, defaultval bool) bool {
	v, err := c.Bool(key)
	if err != nil {
		return defaultval
	}
	return v
}

// Int returns the integer value for a given key.
func (c *Configer) Int(key string) (int, error) {
	return strconv.Atoi(c.getdata(key))
}

// DefaultInt returns the integer value for a given key.
// if err != nil return defaltval
func (c *Configer) DefaultInt(key string, defaultval int) int {
	v, err := c.Int(key)
	if err != nil {
		return defaultval
	}
	return v
}

// Int64 returns the int64 value for a given key.
func (c *Configer) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.getdata(key), 10, 64)
}

// DefaultInt64 returns the int64 value for a given key.
// if err != nil return defaltval
func (c *Configer) DefaultInt64(key string, defaultval int64) int64 {
	v, err := c.Int64(key)
	if err != nil {
		return defaultval
	}
	return v
}

// Float returns the float value for a given key.
func (c *Configer) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.getdata(key), 64)
}

// DefaultFloat returns the float64 value for a given key.
// if err != nil return defaltval
func (c *Configer) DefaultFloat(key string, defaultval float64) float64 {
	v, err := c.Float(key)
	if err != nil {
		return defaultval
	}
	return v
}

// String returns the string value for a given key.
func (c *Configer) Str(key string) string {
	return c.getdata(key)
}

// DefaultString returns the string value for a given key.
// if err != nil return defaltval
func (c *Configer) DefaultString(key string, defaultval string) string {
	v := c.Str(key)
	if v == "" {
		return defaultval
	}
	return v
}

// Strings returns the []string value for a given key.
// Return nil if config value does not exist or is empty.
func (c *Configer) Strings(key string) []string {
	v := c.Str(key)
	if v == "" {
		return nil
	}
	return strings.Split(v, ";")
}

// DefaultStrings returns the []string value for a given key.
// if err != nil return defaltval
func (c *Configer) DefaultStrings(key string, defaultval []string) []string {
	v := c.Strings(key)
	if v == nil {
		return defaultval
	}
	return v
}

func GetConfigData(data [APP_CONF_SIZE]map[string]string, key string) string {
	if len(key) == 0 {
		return ""
	}

	key = strings.ToLower(key)

	if v, ok := data[APP_CONF_FILE][key]; ok {
		return v
	}
	if v, ok := data[APP_CONF_DB][key]; ok {
		return v
	}
	if v, ok := data[APP_CONF_DEFAULT][key]; ok {
		return v
	}
	return ""
}

// section.key or key
func (c *Configer) getdata(key string) string {
	return GetConfigData(c.data, key)
}

// ParseBool returns the boolean value represented by the string.
//
// It accepts 1, 1.0, t, T, TRUE, true, True, YES, yes, Yes,Y, y, ON, on, On,
// 0, 0.0, f, F, FALSE, false, False, NO, no, No, N,n, OFF, off, Off.
// Any other value returns an error.
func ParseBool(val interface{}) (value bool, err error) {
	if val != nil {
		switch v := val.(type) {
		case bool:
			return v, nil
		case string:
			switch v {
			case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "Y", "y", "ON", "on", "On":
				return true, nil
			case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "N", "n", "OFF", "off", "Off":
				return false, nil
			}
		case int8, int32, int64:
			strV := fmt.Sprintf("%s", v)
			if strV == "1" {
				return true, nil
			} else if strV == "0" {
				return false, nil
			}
		case float64:
			if v == 1 {
				return true, nil
			} else if v == 0 {
				return false, nil
			}
		}
		return false, fmt.Errorf("parsing %q: invalid syntax", val)
	}
	return false, fmt.Errorf("parsing <nil>: invalid syntax")
}

func AssignConfig(ac *Configer, ps ...interface{}) error {
	for _, p := range ps {
		assignSingleConfig(ac, p)
	}
	return nil
}

func assignSingleConfig(ac *Configer, p interface{}) {
	pt := reflect.TypeOf(p)
	if pt.Kind() != reflect.Ptr {
		return
	}
	pt = pt.Elem()
	if pt.Kind() != reflect.Struct {
		return
	}
	pv := reflect.ValueOf(p).Elem()

	for i := 0; i < pt.NumField(); i++ {
		pf := pv.Field(i)
		if !pf.CanSet() {
			continue
		}
		name := pt.Field(i).Name
		switch pf.Kind() {
		case reflect.String:
			pf.SetString(ac.DefaultString(name, pf.String()))
		case reflect.Int, reflect.Int64:
			pf.SetInt(int64(ac.DefaultInt64(name, pf.Int())))
		case reflect.Bool:
			pf.SetBool(ac.DefaultBool(name, pf.Bool()))
		case reflect.Struct:
		default:
			//do nothing here
		}
	}
}
