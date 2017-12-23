/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import "fmt"

func (p *GetRequest) Id() string {
	return fmt.Sprintf("%s/%s/%s", p.Endpoint, p.Metric, p.Type)
}

func (p *GetRequest) Csum() string {
	return Md5sum(p.Id())
}

func (p *Item) Id() string {
	return fmt.Sprintf("%s/%s/%s/%s/%d", p.Endpoint, p.Metric, p.Tags, p.Type)
}

func (p *Item) Adjust(now int64) error {
	if p == nil {
		return ErrEmpty
	}

	if len(p.Metric) == 0 || len(p.Endpoint) == 0 {
		return EINVAL
	}

	if len(p.Metric)+len(p.Tags) > 510 {
		return EINVAL
	}

	if p.Ts <= 0 || p.Ts > (now+3600) {
		p.Ts = now
	}

	p.Tags = sortTags(p.Tags)

	return nil
}

func (p *Item) Csum() string {
	return Md5sum(p.Id())
}
