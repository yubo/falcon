/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import "fmt"

func (p *GetRequest) Id() string {
	return fmt.Sprintf("%s/%s/%s/%d", p.Host, p.Name, p.Type, p.Step)
}

func (p *GetRequest) Csum() string {
	return Md5sum(p.Id())
}

func (p *Item) Id() string {
	return fmt.Sprintf("%s/%s/%s/%s/%d", p.Host, p.Name, p.Tags, p.Type, p.Step)
}

func (p *Item) Adjust(now int64) error {
	if p == nil {
		return ErrEmpty
	}

	if len(p.Name) == 0 || len(p.Host) == 0 {
		return EINVAL
	}

	if p.Step <= 0 {
		return EINVAL
	}

	if len(p.Name)+len(p.Tags) > 510 {
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
