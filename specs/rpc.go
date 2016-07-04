package specs

import (
	"fmt"
	"math"
	"strings"
)

// code == 0 => success
// code == 1 => bad request
type RpcResp struct {
	Code int `json:"code"`
}

func (this *RpcResp) String() string {
	return fmt.Sprintf("<Code: %d>", this.Code)
}

type Null struct {
}

/* agent/handoff */
type HandoffResp struct {
	Message string
	Total   int
	Invalid int
}

func (this *HandoffResp) String() string {
	return fmt.Sprintf(
		"<Total=%v, Invalid:%v, Latency=%vms, Message:%s>",
		this.Total,
		this.Invalid,
		this.Message,
	)
}

type MetaData struct {
	Host string  `json:"host"`
	K    string  `json:"k"`
	V    float64 `json:"v"`
	Ts   int64   `json:"ts"`
	Step int64   `json:"step"`
	Type string  `json:"type"`
	Tags string  `json:"tags"`
}

func (t *MetaData) String() string {
	return fmt.Sprintf("<MetaData host:%s, metric:%s, Timestamp:%d, Step:%d, Value:%f, Tags:%v>",
		t.Host, t.K, t.Ts, t.Step, t.V, t.Tags)
}

func (p *MetaData) Key() string {
	if p.Tags == "" {
		return fmt.Sprintf("%s/%s", p.Host, p.K)
	}
	return fmt.Sprintf("%s/%s/%s", p.Host, p.K, p.Tags)
}

func (p *MetaData) Rrd() (*RrdItem, error) {
	e := &RrdItem{}

	e.Endpoint = p.Host
	e.Metric = p.K
	e.Tags = p.Tags
	e.Timestamp = p.Ts
	e.Value = p.V
	e.Step = int(p.Step)
	if e.Step < MIN_STEP {
		e.Step = MIN_STEP
	}
	e.Heartbeat = e.Step * 2

	if p.Type == GAUGE {
		e.DsType = p.Type
		e.Min = "U"
		e.Max = "U"
	} else if p.Type == COUNTER {
		e.DsType = DERIVE
		e.Min = "0"
		e.Max = "U"
	} else if p.Type == DERIVE {
		e.DsType = DERIVE
		e.Min = "0"
		e.Max = "U"
	} else {
		return e, fmt.Errorf("not_supported_counter_type")
	}

	e.Timestamp = e.Timestamp - e.Timestamp%int64(e.Step)

	return e, nil
}

func (p *MetaData) Tsdb() *TsdbItem {
	t := TsdbItem{Tags: make(map[string]string)}

	if p.Tags != "" {
		tags := strings.Split(p.Tags, ",")
		for _, tag := range tags {
			kv := strings.SplitN(tag, "=", 2)
			if len(kv) == 2 {
				t.Tags[kv[0]] = kv[1]
			}
		}
	}
	t.Tags["host"] = p.Host
	t.Metric = p.K
	t.Timestamp = p.Ts
	t.Value = p.V
	return &t
}

type TsdbItem struct {
	Metric    string            `json:"metric"`
	Tags      map[string]string `json:"tags"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
}

func (this *TsdbItem) String() string {
	return fmt.Sprintf(
		"<Metric:%s, Tags:%v, Value:%v, TS:%d>",
		this.Metric,
		this.Tags,
		this.Value,
		this.Timestamp,
	)
}

func (this *TsdbItem) TsdbString() (s string) {
	s = fmt.Sprintf("put %s %d %.3f ", this.Metric, this.Timestamp, this.Value)

	for k, v := range this.Tags {
		key := strings.ToLower(strings.Replace(k, " ", "_", -1))
		value := strings.Replace(v, " ", "_", -1)
		s += key + "=" + value + " "
	}

	return s
}

/* handoff/storage */
// DsType 即RRD中的Datasource的类型：GAUGE|COUNTER|DERIVE
type RrdItem struct {
	Endpoint  string  `json:"endpoint"`
	Metric    string  `json:"metric"`
	Tags      string  `json:"tags"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
	DsType    string  `json:"dstype"`
	Step      int     `json:"step"`
	Heartbeat int     `json:"heartbeat"`
	Min       string  `json:"min"`
	Max       string  `json:"max"`
}

func (this *RrdItem) String() string {
	return fmt.Sprintf(
		"<Endpoint:%s, Metric:%s, Tags:%v, Value:%v, TS:%d %v DsType:%s, Step:%d, Heartbeat:%d, Min:%s, Max:%s>",
		this.Endpoint,
		this.Metric,
		this.Tags,
		this.Value,
		this.Timestamp,
		fmtTs(this.Timestamp),
		this.DsType,
		this.Step,
		this.Heartbeat,
		this.Min,
		this.Max,
	)
}

func (p *RrdItem) Csum() string {
	return md5sum(p.Id())
}

func (p *RrdItem) Id() string {
	return fmt.Sprintf("%s/%s/%s/%s/%d", p.Endpoint, p.Metric, p.Tags, p.DsType, p.Step)
}

// ConsolFun 是RRD中的概念，比如：MIN|MAX|AVERAGE
type RrdQuery struct {
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	ConsolFun string `json:"consolFuc"`
	Endpoint  string `json:"endpoint"`
	Counter   string `json:"counter"`
	DsType    string `json:"dstype"`
	Step      int    `json:"step"`
}

func (p *RrdQuery) Csum() string {
	return md5sum(p.Id())
}

func (p *RrdQuery) Id() string {
	return fmt.Sprintf("%s/%s/%s/%d", p.Endpoint, p.Counter, p.DsType, p.Step)
}

type RrdResp struct {
	Endpoint string     `json:"endpoint"`
	Counter  string     `json:"counter"`
	DsType   string     `json:"dstype"`
	Step     int        `json:"step"`
	Values   []*RRDData `json:"Values"` //大写为了兼容已经再用这个api的用户
}

func (p *RrdResp) Csum() string {
	return md5sum(p.Id())
}

func (p *RrdResp) Id() string {
	return fmt.Sprintf("%s/%s/%s/%d", p.Endpoint, p.Counter, p.DsType, p.Step)
}

type RrdQueryCsum struct {
	Csum      string `json:"csum"`
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	ConsolFun string `json:"consolFuc"`
}

type RrdRespCsum struct {
	Values []*RRDData `json:"values"`
}

type JsonFloat float64

func (v JsonFloat) MarshalJSON() ([]byte, error) {
	f := float64(v)
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return []byte("null"), nil
	} else {
		return []byte(fmt.Sprintf("%f", f)), nil
	}
}

type RRDData struct {
	Timestamp int64     `json:"timestamp"`
	Value     JsonFloat `json:"value"`
}

func (this *RRDData) String() string {
	return fmt.Sprintf(
		"<RRDData:Value:%v TS:%d %v>",
		this.Value,
		this.Timestamp,
		fmtTs(this.Timestamp),
	)
}

type File struct {
	Filename string
	Data     []byte
}
