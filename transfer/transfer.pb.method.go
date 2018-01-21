/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/service"
)

func adjustKey(key []byte) ([]byte, error) {
	endpoint, metric, tags, typ, err := falcon.KeyAttr(key)
	if err != nil {
		return key, err
	}

	if len(tags) > 0 {
		tags_ := strings.Split(tags, ",")
		sort.Strings(tags_)
		tags = strings.Join(tags_, ",")
		return []byte(fmt.Sprintf("%s/%s/%s/%s", endpoint, metric, tags, typ)), nil
	}
	return key, nil
}

// call before use item.Key
func (p *DataPoint) Adjust() (*service.DataPoint, error) {
	key, err := adjustKey(p.Key)
	if err != nil {
		return nil, err
	}

	return &service.DataPoint{
		Key:   &service.Key{Key: key, ShardId: int32(scheduler(key))},
		Value: p.Value,
	}, nil
}

func (p *GetRequest) Adjust() (map[int]*service.GetRequest, error) {
	ret := make(map[int]*service.GetRequest)

	for _, key := range p.Keys {
		shardId := scheduler(key)
		if req, ok := ret[shardId]; ok {
			req.Keys = append(req.Keys, &service.Key{
				Key:     key,
				ShardId: int32(shardId),
			})
		} else {
			ret[shardId] = &service.GetRequest{
				Keys: []*service.Key{&service.Key{
					Key:     key,
					ShardId: int32(shardId),
				}},
				Start:     p.Start,
				End:       p.End,
				ConsolFun: p.ConsolFun,
			}
		}

	}
	return ret, nil
}
