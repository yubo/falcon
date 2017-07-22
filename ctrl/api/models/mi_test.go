/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package models

import (
	"fmt"
	"testing"
)

func testMi(t *testing.T) {
	miSyncer.ctx, _ = miFetchHosts("/home/yubo/tmp/hostinfos")
	h, err := miGetTagHost("cop=xiaomi,owt=dba,pdl=huyu,service=game-stat", "",
		true, 0, 0)
	fmt.Println(len(h), h, err)

}
