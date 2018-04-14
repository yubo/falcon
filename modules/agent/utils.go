/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"strings"

	"github.com/yubo/falcon/lib/core"
)

func counterAttr(key string) (string, string, string, error) {
	var err error
	s := strings.Split(key, "/")
	if len(s) != 3 {
		err = core.EINVAL
	}

	return string(s[0]), string(s[1]), string(s[2]), err
}

func counterFmtOk(counter string) bool {
	_, _, _, err := counterAttr(counter)
	return err != nil
}
