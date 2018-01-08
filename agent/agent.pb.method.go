/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"fmt"

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/service"
)

func (p *Item) toServiceItem(endpoint []byte) (*service.Item, error) {
	if len(endpoint) == 0 || len(p.Metric) == 0 {
		return nil, falcon.EINVAL
	}

	if _, ok := service.ItemType_value[string(p.Type)]; !ok {
		return nil, falcon.EINVAL
	}

	return &service.Item{
		Key:       []byte(fmt.Sprintf("%s/%s/%s/%s", endpoint, p.Metric, p.Tags, p.Type)),
		Value:     p.Value,
		Timestamp: p.Timestamp,
	}, nil
}
