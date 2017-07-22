/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"errors"
	"math"
	"net"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcModule struct {
	enable  bool
	ctx     context.Context
	cancel  context.CancelFunc
	address string
	b       *Backend
}

func (p *GrpcModule) queryGetCacheEntry(in *falcon.GetRequest) (*cacheEntry, error) {

	e := p.b.cache.get(in.Csum())
	if e == nil {
		return nil, falcon.ErrNoExits
	}

	in.Start = in.Start - in.Start%int64(in.Step)
	in.End = in.End - in.End%int64(in.Step) + int64(in.Step)
	if in.End-in.Start-int64(in.Step) < 1 {
		return nil, falcon.ErrParam
	}
	return e, nil
}

func (p *GrpcModule) queryGetData(in *falcon.GetRequest, resp *falcon.GetResponse,
	e *cacheEntry) (rrds, caches []*falcon.RRDData, err error) {

	flag := atomic.LoadUint32(&e.flag)
	caches, _ = e._getData(e.commitId, e.dataId)

	if !p.b.Conf.Migrate.Disabled && flag&RRD_F_MISS != 0 {
		node, _ := p.b.storageMigrateConsistent.Get(in.Id())
		done := make(chan error, 1)
		res := []*falcon.RRDData{}
		p.b.storageNetTaskCh[node] <- &netTask{
			Method: NET_TASK_M_QUERY,
			Done:   done,
			Args:   in,
			Reply:  res,
		}
		<-done
		// fetch data from remote
		rrds = res
	} else {
		// read data from local rrd file
		rrds, _ = p.b.taskRrdFetch(e.hashkey, in.ConsolFun,
			in.Start, in.End, int(e.step))
	}

	// larger than rra1point range, skip merge
	now := p.b.timeNow()
	if in.Start < now-now%int64(e.step)-RRA1PointCnt*int64(e.step) {
		resp.Vs = rrds
		return nil, nil, errors.New("skip merge")
	}

	// no cached caches, do not merge
	if len(caches) < 1 {
		resp.Vs = rrds
		return nil, nil, errors.New("no caches")
	}
	return rrds, caches, nil
}

func queryPruneCache(items []*falcon.RRDData, e *cacheEntry,
	start, end int64) (ret []*falcon.RRDData) {

	// prune cached items
	var val float64

	ts := items[0].Ts
	n := len(items)

	last := items[n-1].Ts
	i := 0
	if e.typ == falcon.ItemType_DERIVE || e.typ == falcon.ItemType_COUNTER {
		for ts < last {
			if i < n-1 && ts == items[i].Ts &&
				ts == items[i+1].Ts-int64(e.step) {
				val = (items[i+1].V - items[i].V) / float64(e.step)
				if val < 0 {
					val = math.NaN()
				}
				i++
			} else {
				// missing
				val = math.NaN()
			}

			if ts >= start && ts <= end {
				ret = append(ret,
					&falcon.RRDData{Ts: ts, V: val})
			}
			ts = ts + int64(e.step)
		}
	} else if e.typ == falcon.ItemType_GAUGE {
		for ts <= last {
			if i < n && ts == items[i].Ts {
				val = items[i].V
				i++
			} else {
				// missing
				val = math.NaN()
			}

			if ts >= start && ts <= end {
				ret = append(ret,
					&falcon.RRDData{Ts: ts, V: val})
			}
			ts = ts + int64(e.step)
		}
	}
	return ret
}

/*
 * a older than b
 * c = a <- b
 */
func queryMergeData(a, b []*falcon.RRDData, start,
	end, step int64) []*falcon.RRDData {

	// do merging
	c := make([]*falcon.RRDData, 0)
	if len(a) > 0 {
		for _, v := range a {
			if v.Ts >= start &&
				v.Ts <= end {
				//rrdtool返回的数据,时间戳是连续的、不会有跳点的情况
				c = append(c, v)
			}
		}
	}

	bl := len(b)
	if bl > 0 {
		cl := len(c)
		lastTs := b[0].Ts

		// find junction
		i := 0
		for i = cl - 1; i >= 0; i-- {
			if c[i].Ts < b[0].Ts {
				lastTs = c[i].Ts
				break
			}
		}

		// fix missing
		for ts := lastTs + step; ts < b[0].Ts; ts += step {
			c = append(c, &falcon.RRDData{Ts: ts,
				V: math.NaN()})
		}

		// merge cached items to result
		i += 1
		for j := 0; j < bl; j++ {
			if i < cl {
				if !math.IsNaN(float64(b[j].V)) {
					c[i] = b[j]
				}
			} else {
				c = append(c, b[j])
			}
			i++
		}
	}
	return c
}

func queryPruneRet(a []*falcon.RRDData,
	start, end, step int64) []*falcon.RRDData {

	// prune result
	n := int((end - start) / step)
	ret := make([]*falcon.RRDData, n)
	j := 0
	ts := start
	al := len(a)

	for i := 0; i < n; i++ {
		if j < al && ts == a[j].Ts {
			ret[i] = a[j]
			j++
		} else {
			ret[i] = &falcon.RRDData{Ts: ts,
				V: math.NaN()}
		}
		ts += step
	}
	return ret
}

func (p *GrpcModule) Get(ctx context.Context,
	in *falcon.GetRequest) (resp *falcon.GetResponse, err error) {
	var (
		e    *cacheEntry
		rrds []*falcon.RRDData
		ret  []*falcon.RRDData
	)

	statsInc(ST_RPC_SERV_QUERY, 1)

	resp = &falcon.GetResponse{
		Vs:   []*falcon.RRDData{},
		Host: in.Host,
		Name: in.Name,
		Type: in.Type,
		Step: in.Step,
	}

	if e, err = p.queryGetCacheEntry(in); err != nil {
		return
	}

	if rrds, ret, err = p.queryGetData(in, resp, e); err != nil {
		statsInc(ST_RPC_SERV_QUERY_ITEM, len(resp.Vs))
		return
	}

	ret = queryPruneCache(ret, e, in.Start, in.End)
	ret = queryMergeData(rrds, ret, in.Start, in.End, int64(e.step))
	resp.Vs = queryPruneRet(ret, in.Start, in.End, int64(e.step))

	statsInc(ST_RPC_SERV_QUERY_ITEM, len(resp.Vs))

	return
}

func (p *GrpcModule) GetRrd(ctx context.Context,
	in *falcon.GetRrdRequest) (*falcon.GetRrdResponse, error) {

	statsInc(ST_RPC_SERV_GETRRD, 1)

	key := string(in.Key)
	e := p.b.cache.get(key)
	if e != nil {
		e.commit(p.b)
	}

	data, err := p.b.taskFileRead(key)
	if err != nil {
		statsInc(ST_RPC_SERV_GETRRD_ERR, 1)
	}

	return &falcon.GetRrdResponse{File: &falcon.File{Data: data}}, nil
}

func (p *GrpcModule) Update(ctx context.Context,
	in *falcon.UpdateRequest) (*falcon.UpdateResponse, error) {

	total, errors := p.b.handleItems(in.Items)
	return &falcon.UpdateResponse{Total: int32(total), Errors: int32(errors)}, nil
}

func (p *GrpcModule) prestart(backend *Backend) error {
	p.enable, _ = backend.Conf.Configer.Bool(C_GRPC_ENABLE)
	p.address = backend.Conf.Configer.Str(C_GRPC_ADDR)
	p.b = backend
	return nil
}

func (p *GrpcModule) start(backend *Backend) error {

	if !p.enable {
		glog.Info(MODULE_NAME + "grpc.Start not enabled")
		return nil
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	ln, err := net.Listen(falcon.ParseAddr(p.address))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	RegisterBackendServer(server, &GrpcModule{})

	// Register reflection service on gRPC server.
	reflection.Register(server)
	go func() {
		if err := server.Serve(ln); err != nil {
			p.cancel()
		}
	}()

	go func() {
		<-p.ctx.Done()
		server.Stop()
	}()

	return nil
}

func (p *GrpcModule) stop(backend *Backend) error {
	if !p.enable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *GrpcModule) reload(backend *Backend) error {
	if p.enable {
		p.stop(backend)
		time.Sleep(time.Second)
	}
	p.prestart(backend)
	return p.start(backend)
}
