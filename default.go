/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package main

const (
	VERSION      = "0.0.3"
	DEFAULT_STEP = 60 //s
	SHARD_NUM    = 10
	//MIN_STEP     = 30 //s
	REPLICAS = 500
	//GAUGE       = "GAUGE"
	//DERIVE      = "DERIVE"
	//COUNTER     = "COUNTER"
	//MODULE_NAME = "\x1B[37m[FALCON]\x1B[0m"
)

const (
	BaseConf = `
agent:
  disable: true
  burst_size: 16	# client put burst size to remote service
  interval: 5
  workerProcesses: 3
  iface_prefix:
    - eth
    - em

service:
  disable: true
  call_timeout: 5000
  idx: true
  db_max_idle: 4
  db_max_conn: 4
  conf_interval: 600
  tsdb_bucket_num: 13
  tsdb_bucket_size: 7200
  tsdb_dir: /tmp/tsdb

sys:
  log_level: 4
  log_file: stdout
  pid_file: /var/run/falcon.pid

`
)
