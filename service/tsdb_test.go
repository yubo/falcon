package service

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/huangaz/tsdb/lib/testUtil"
	fconfig "github.com/yubo/falcon/config"
	"github.com/yubo/falcon/service/config"
)

var (
	s    *Service
	tsdb *TsdbModule
)

func test_init() {
	s = &Service{
		Conf: &config.Service{
			Name: "cacheApp",
		},
	}
	s.Conf.Configer.Set(fconfig.APP_CONF_FILE, map[string]string{
		"shardIds": "1",
	})

	tsdb = &TsdbModule{}
	tsdb.prestart(s)
}

func TestPut(t *testing.T) {
	test_init()
	tsdb.start(s)
	defer tsdb.stop(s)

	num := 10
	putReq := dataGenerator(1, num)

	time.Sleep(1000 * time.Millisecond)

	putRes, err := tsdb.put(putReq)
	if err != nil {
		t.Fatal(err)
	}
	if putRes.N != int32(num) {
		t.Fatalf("Put failed! Want %d,put %d", num, putRes.N)
	}

	time.Sleep(1000 * time.Millisecond)

	getReq := &GetRequest{
		Start:   0,
		End:     int64(60 * (num + 1)),
		ShardId: 1,
		Key:     putReq.Items[0].Key,
	}
	getRes, err := tsdb.get(getReq)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(putReq.PrintForDebug())
	fmt.Println(getRes.PrintForDebug())

	if string(getRes.Key) != string(putReq.Items[0].Key) {
		t.Fatal("wrong result")
	}

	if len(putReq.Items) != len(getRes.Dps) {
		t.Fatalf("Length of putReq(%d) and getRes(%d) not equal!", len(putReq.Items), len(getRes.Dps))
	}

	for i, item := range putReq.Items {
		if item.Value != getRes.Dps[i].Value || item.Timestamp != getRes.Dps[i].Timestamp {
			t.Fatal("wrong result")
		}
	}
}

func dataGenerator(numOfKeys, num int) *PutRequest {
	req := &PutRequest{}
	req.Items = make([]*Item, num*numOfKeys)
	index := 0

	for i := 0; i < numOfKeys; i++ {
		testKey := []byte(testUtil.RandStr(10))
		testTime := 0

		for j := 0; j < num; j++ {
			testTime += (55 + rand.Intn(10))
			newItem := &Item{
				Key:       testKey,
				ShardId:   int32(i + 1),
				Timestamp: int64(testTime),
				Value:     float64(100 + rand.Intn(50)),
			}
			req.Items[index] = newItem
			index++
		}
	}
	return req
}
