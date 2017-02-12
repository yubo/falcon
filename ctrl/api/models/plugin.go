/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

type hookfunc func() error

var (
	hooks = make([]hookfunc, 0) //hook function slice to store the hookfunc
)

func RegisterPlugin(hf hookfunc) {
	hooks = append(hooks, hf)
}

func PluginStart() {
	for _, hk := range hooks {
		if err := hk(); err != nil {
			panic(err)
		}
	}
}
