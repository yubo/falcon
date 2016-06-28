package storage

import (
	"io/ioutil"
	"os"
	"sync/atomic"
	"unsafe"

	"github.com/golang/glog"
	"github.com/yubo/falcon/hcl"
)

func Config() *Options {
	return (*Options)(atomic.LoadPointer(&configPtr))
}

func applyConfigFile(opts *Options, filePath string) error {
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

	return nil
}

func parse(config string) {

	_, err := os.Stat(config)
	if !os.IsNotExist(err) {
		if err := applyConfigFile(&configOpts, config); err != nil {
			glog.Errorln(err)
			os.Exit(2)
		}
	}

	if configOpts.Migrate.Enable && len(configOpts.Migrate.Cluster) == 0 {
		configOpts.Migrate.Enable = false
	}

	// set config
	atomic.StorePointer(&configPtr, unsafe.Pointer(&configOpts))

	glog.V(3).Infof("ParseConfig ok, file %s", config)
	configFile = config
}
