/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import "fmt"

func (p *GetRequest) Key() string {
	return fmt.Sprintf("%s/%s/%s/%d", p.Endpoint, p.Metric, p.Tags, p.Type)
}

func (p *GetRequest) Csum() string {
	return Md5sum(p.Key())
}

func (p *GetRequest) Sum64() uint64 {
	return Sum64(p.Key())
}

func (p *GetRequest) Sum32() uint32 {
	return Sum32(p.Key())
}

func (p *Item) Key() string {
	return fmt.Sprintf("%s/%s/%s/%d", p.Endpoint, p.Metric, p.Tags, p.Type)
}

func (p *Item) Sum64() uint64 {
	return Sum64(p.Key())
}

func (p *Item) Sum32() uint32 {
	return Sum32(p.Key())
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

	if p.Timestamp <= 0 || p.Timestamp > (now+3600) {
		p.Timestamp = now
	}

	p.Tags = sortTags(p.Tags)

	return nil
}

func (p *Item) Csum() string {
	return Md5sum(p.Key())
}
