/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

type TagHostApiGet struct {
	Id            int64  `json:"id"`
	TagId         int64  `json:"tag_id"`
	HostId        int64  `json:"host_id"`
	TagName       string `json:"tag_name"`
	HostName      string `json:"host_name"`
	Pause         int64  `json:"pause"`
	MaintainBegin int64  `json:"maintain_begin"`
	MaintainEnd   int64  `json:"maintain_end"`
}
