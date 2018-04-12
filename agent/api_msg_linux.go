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

	"github.com/yubo/falcon/lib/ipc"
	"github.com/yubo/falcon/lib/tsdb"
	"github.com/yubo/falcon/transfer"
)

func init() {
	RegisterModule(&ApiMsgModule{})
}

const (
	IPC_MQ_MODE = 0666
	IPC_MQ_KEY  = 0x1234
	IPC_MQ_TYPE = 1
)

type ApiMsgModule struct {
	disable bool

	ctx    context.Context
	cancel context.CancelFunc
}

func (p *ApiMsgModule) prestart(agent *Agent) error {
	p.disable = !agent.Conf.IpcEnable
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
		msqid := ipc.MsgGet(IPC_MQ_KEY, IPC_MQ_MODE)
		for {
			select {
			case <-p.ctx.Done():
				return
			default:
				msg, err := ipc.MsgRcv(msqid, IPC_MQ_TYPE, 0)
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
	p.disable = !agent.Conf.IpcEnable
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
