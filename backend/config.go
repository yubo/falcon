/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import "github.com/yubo/falcon/specs"

const (
	CTRL_STEP = 360
)

var (
	DefaultBackend = Backend{
		Params: specs.ModuleParams{
			Debug:       0,
			ConnTimeout: 1000,
			CallTimeout: 5000,
			Concurrency: 2,
			Disabled:    false,
			Http:        true,
			Rpc:         true,
			HttpAddr:    "0.0.0.0:7021",
			RpcAddr:     "0.0.0.0:7020",
			CtrlAddr:    "",
		},
		Migrate: specs.Migrate{
			Disabled:  false,
			Upstreams: map[string]string{},
		},
		Idx:             true,
		IdxInterval:     30,
		IdxFullInterval: 86400,
		Dsn:             "falcon:1234@tcp(127.0.0.1:3306)/falcon?loc=Local&parseTime=true",
		DbMaxIdle:       4,
		ShmMagic:        0x80386,
		ShmKey:          0x7020,
		ShmSize:         256 * (1 << 20), // 256m
		Storage: Storage{
			Type: "rrd",
		},
	}
)

/*
func applyConfigFile(opts *BackendOpts, filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}

	fileString := []byte{}
	glog.V(3).Infof("Loading config file at: %s", filePath)

	fileString, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := hcl.Decode(opts, string(fileString)); err != nil {
		return err
	}

	glog.V(3).Infof("config options:\n%s", opts)
	return nil
}

func parse(config *BackendOpts, filename string) {

	if err := applyConfigFile(config, filename); err != nil {
		glog.Errorln(err)
		os.Exit(2)
	}

	if config.Migrate.Enable && len(config.Migrate.Upstreams) == 0 {
		config.Migrate.Enable = false
	}

	// set config
	//atomic.StorePointer(&configPtr, unsafe.Pointer(&configOpts))

	glog.V(3).Infof("ParseConfig ok, file %s", filename)
	glog.V(3).Infof("cache_size %d", CACHE_SIZE)
	appConfigfile = filename
}
*/
