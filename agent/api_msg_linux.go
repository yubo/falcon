/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/yubo/falcon/lib/tsdb"
	"github.com/yubo/falcon/transfer"
)

type ApiMsgModule struct {
	disable bool

	ctx    context.Context
	cancel context.CancelFunc
}

func (p *ApiMsgModule) prestart(agent *Agent) error {
	p.disable = !agent.Conf.Configer.DefaultBool(C_IPC_ENABLE, false)
	return nil
}
func (p *ApiMsgModule) start(agent *Agent) error {
	if p.disable {
		return nil
	}
	p.ctx, p.cancel = context.WithCancel(context.Background())
	endpoint := agent.Conf.Host
	ch := agent.PutChan

	go func() {
		msqid := MsgGet(IPC_MQ_KEY, IPC_MQ_MODE)
		for {
			select {
			case <-p.ctx.Done():
				return
			default:
				msg, err := MsgRcv(msqid, IPC_MQ_TYPE, 0)
				if err != nil {
					continue
				}

				if req := msgToReq(endpoint, msg); req != nil {
					ch <- req
				}

			}
		}
	}()
	return nil
}

func (p *ApiMsgModule) stop(agent *Agent) error {
	if p.disable {
		return nil
	}
	p.cancel()
	return nil
}
func (p *ApiMsgModule) reload(agent *Agent) error {
	p.stop(agent)
	p.disable = !agent.Conf.Configer.DefaultBool(C_IPC_ENABLE, false)
	return p.start(agent)
}

func msgToReq(endpoint, msg string) *PutRequest {
	var (
		counter string
		value   float64
	)
	n, err := fmt.Sscanf(msg, "%s %f", &counter, &value)
	if err != nil || n != 2 || !counterFmtOk(counter) {
		return nil
	}

	return &PutRequest{
		Dps: []*transfer.DataPoint{
			&transfer.DataPoint{
				Key: []byte(endpoint + "/" + counter),
				Value: &tsdb.TimeValuePair{
					Timestamp: time.Now().Unix(),
					Value:     value,
				},
			},
		},
	}
}
