/*
 * Copyright 2018 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

/*
 * the configer is not thread safe,
 * make sure not use it after call process.Start()
 */

package core

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/ghodss/yaml"
)

type Configer struct {
	configFile string
	valueFiles []string
	base       map[string]interface{}
	data       map[string]interface{}
}

func ToConfiger(in interface{}) *Configer {
	if data, ok := in.(map[string]interface{}); ok {
		return &Configer{data: data}
	}
	return &Configer{data: map[string]interface{}{}}
}

func NewConfiger(configFile, baseConf string, valueFiles []string) (conf *Configer, err error) {
	conf = &Configer{
		base: map[string]interface{}{},
		data: map[string]interface{}{},
	}

	for _, file := range valueFiles {
		file, err = filepath.Abs(file)
		if err != nil {
			return
		}
		conf.valueFiles = append(valueFiles, file)
	}

	conf.configFile, err = filepath.Abs(configFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal([]byte(baseConf), &conf.base)
	if err != nil {
		return
	}

	return
}

func (p *Configer) GetConfiger(path string) *Configer {
	return ToConfiger(p.GetRaw(path))
}

func (p *Configer) Parse() error {
	values := map[string]interface{}{}

	// parse values file to values
	for _, file := range p.valueFiles {
		b, err := ParseTemplateFile(values, file)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %s", file, err)
		}

		values, err = yaml2Values(values, b)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %s", file, err)
		}
	}

	// parse config file with values
	b, err := ParseTemplateFile(
		map[string]interface{}{"values": values},
		p.configFile)
	if err != nil {
		return err
	}
	p.data, err = yaml2Values(p.base, b)
	return err
}

func (p *Configer) GetRaw(path string) interface{} {
	v, err := Values(p.data).PathValue(path)
	if err != nil {
		return nil
	}
	return v
}

func (p *Configer) GetStr(path string) string {
	v, err := Values(p.data).PathValue(path)
	if err != nil {
		return ""
	}

	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

func (p *Configer) GetBool(path string) (bool, error) {
	v, err := Values(p.data).PathValue(path)
	if err != nil {
		return false, err
	}

	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("%v is not bool", path)
	}
	return b, nil
}

func (p *Configer) GetBoolDef(path string, def bool) bool {
	v, err := p.GetBool(path)
	if err != nil {
		return def
	}
	return v
}

func (p *Configer) GetFloat64(path string) (float64, error) {
	v, err := Values(p.data).PathValue(path)
	if err != nil {
		return 0, err
	}

	f, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("%v is not number", path)
	}

	return f, nil
}

func (p *Configer) GetFloat64Def(path string, def float64) float64 {
	v, err := p.GetFloat64(path)
	if err != nil {
		return def
	}

	return v
}

func (p *Configer) GetInt64(path string) (int64, error) {
	v, err := p.GetFloat64(path)
	if err != nil {
		return 0, err
	}

	return int64(v), nil
}

func (p *Configer) GetInt64Def(path string, def int64) int64 {
	v, err := p.GetInt64(path)
	if err != nil {
		return def
	}
	return v
}

func (p *Configer) GetInt(path string) (int, error) {
	v, err := p.GetFloat64(path)
	if err != nil {
		return 0, err
	}

	return int(v), nil
}

func (p *Configer) GetIntDef(path string, def int) int {
	v, err := p.GetInt(path)
	if err != nil {
		return def
	}
	return v
}

func (p *Configer) Read(path string, o interface{}) error {
	v, err := Values(p.data).PathValue(path)
	if err != nil {
		return ErrNoExits
	}

	if o == nil {
		return nil
	}

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, o)
}

func (p *Configer) Unmarshal(o interface{}) error {
	data, err := json.Marshal(p.data)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, o)
}

// just for dump
func _yaml(v map[string]interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func (p *Configer) String() string {
	return _yaml(p.data)
	/*
		return _yaml(map[string]interface{}{
			"tpl":    p.tpl,
			"base":   _yaml(p.base),
			"values": _yaml(p.values),
			"data":   _yaml(p.data),
		})
	*/
}

func yaml2Values(dest map[string]interface{}, bytes []byte) (map[string]interface{}, error) {
	currentMap := map[string]interface{}{}
	if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
		return dest, err
	}
	dest = mergeValues(dest, currentMap)

	return dest, nil
}

// Merges source and destination map, preferring values from the source map
func mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = nextMap
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}
