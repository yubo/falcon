/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

const (
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60
)

var (
	DefaultAgent = Agent{
		Debug:    0,
		Name:     "Agent Module",
		Host:     "",
		Http:     true,
		HttpAddr: "127.0.0.1:1988",
		Rpc:      true,
		RpcAddr:  "127.0.0.1:1989",
		IfPre:    []string{"eth", "em"},
		Interval: 60,
		Handoff: Handoff{
			Batch:       16,
			ConnTimeout: 1000,
			CallTimeout: 5000,
			Upstreams:   []string{},
		},
	}
)

/*
func applyConfigFile(opts *AgentOpts, filePath string) error {
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

func parse(config *AgentOpts, filename string) {

	if err := applyConfigFile(config, filename); err != nil {
		glog.Errorln(err)
		os.Exit(2)
	}
	if config.Host == "" {
		config.Host, _ = os.Hostname()
	}

	glog.V(3).Infof("ParseConfig ok, file %s", filename)
	appConfigFile = filename
}
*/
