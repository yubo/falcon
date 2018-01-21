/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yubo/falcon"
)

func (p *Key) Sum64() uint64 {
	return falcon.Sum64(p.Key)
}

func (p *Key) Sum32() uint32 {
	return falcon.Sum32(p.Key)
}

func (p *Key) Csum() string {
	return falcon.Md5sum(p.Key)
}

// call before use item.Key
func (p *Key) Adjust() error {
	endpoint, metric, tags, typ, err := p.Attr()
	if err != nil {
		return err
	}

	if len(tags) > 0 {
		tags_ := strings.Split(tags, ",")
		sort.Strings(tags_)
		tags = strings.Join(tags_, ",")
		p.Key = []byte(fmt.Sprintf("%s/%s/%s/%s", endpoint, metric, tags, typ))
	}

	return nil
}

func (p *Key) Attr() (string, string, string, string, error) {
	var err error
	s := strings.Split(string(p.Key), "/")
	if len(s) != 4 {
		err = falcon.EINVAL
	}

	if _, ok := ItemType_value[s[3]]; !ok {
		err = falcon.EINVAL
	}
	return s[0], s[1], s[2], s[3], err
}
