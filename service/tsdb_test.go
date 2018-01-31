package service

import (
	"math/rand"
	"testing"
	"time"

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/lib/tsdb"
	"github.com/yubo/falcon/service/config"
)

var (
	s *Service
)

func test_init() *TsdbModule {
	s = &Service{
		Conf: &config.Service{
			Name: "cacheApp",
		},
	}
	s.Conf.Configer.Set(falcon.APP_CONF_FILE, map[string]string{
		"shardIds": "1,2,3,4",
	})

	tm := &TsdbModule{}
	tm.prestart(s)

	return tm
}

/*
func TestPut(t *testing.T) {
	tm := test_init()
	tm.start(s)
	defer tm.stop(s)

	num := 10
	putReq := dataGenerator(0, 1, num,1)

	time.Sleep(100 * time.Millisecond)

	putRes, err := tm.put(putReq)
	if err != nil {
		t.Fatal(err)
	}
	if putRes.N != int32(num) {
		t.Fatalf("Put failed! Want %d,put %d", num, putRes.N)
	}

	time.Sleep(100 * time.Millisecond)

	getReq := &GetRequest{
		Start: 0,
		End:   int64(60 * num),
		Keys: []*tsdb.Key{
			&tsdb.Key{
				ShardId: 1,
				Key:     putReq.Data[0].Key.Key,
			},
		},
	}
	getRes, err := tm.get(getReq)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(putReq.PrintForDebug())
	fmt.Println(getRes.PrintForDebug())

	if string(getRes.Data[0].Key.Key) != string(putReq.Data[0].Key.Key) {
		t.Fatal("wrong result")
	}

	if len(putReq.Data) != len(getRes.Data[0].Values) {
		t.Fatalf("Length of putReq(%d) and getRes(%d) not equal!", len(putReq.Data), len(getRes.Data))
	}

	for i, dp := range putReq.Data {
		if dp.Value.Value != getRes.Data[0].Values[i].Value ||
			dp.Value.Timestamp != getRes.Data[0].Values[i].Timestamp {
			t.Fatal("wrong result")
		}
	}
}
*/

func TestPutAndRecover(t *testing.T) {
	tm := test_init()
	tm.start(s)

	num := 1000
	numOfKeys := 10
	beginTimestamp := time.Now().Unix() - int64(num*60)
	putReq := dataGenerator(beginTimestamp, numOfKeys, num, 1)

	time.Sleep(100 * time.Millisecond)

	putRes, err := tm.put(putReq)
	if err != nil {
		t.Fatal(err)
	}
	if putRes.N != int32(num*numOfKeys) {
		t.Fatalf("Put failed! Want %d,put %d", num*numOfKeys, putRes.N)
	}

	time.Sleep(1000 * time.Millisecond)

	bucketToFinalize := tm.buckets[1].Bucket(beginTimestamp)
	n, err := tm.buckets[1].FinalizeBuckets(bucketToFinalize)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatal("finalize failed")
	}
	// log.Println("lastFinalizedBucket_: ", tm.buckets[1].GetLastFinalizedBucket())

	tm.stop(s)

	// recover data from disk
	t2 := test_init()
	t2.start(s)
	time.Sleep(100 * time.Millisecond)

	getReq := &GetRequest{
		Start: 0,
		End:   time.Now().Unix(),
		Keys: []*tsdb.Key{
			&tsdb.Key{
				ShardId: 1,
				Key:     putReq.Data[0].Key.Key,
			},
		},
	}
	getRes, err := t2.get(getReq)
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Println(putReq.PrintForDebug())
	// fmt.Println(getRes.PrintForDebug())

	if string(getRes.Data[0].Key.Key) != string(putReq.Data[0].Key.Key) {
		t.Fatal("wrong result")
	}

	if num != len(getRes.Data[0].Values) {
		t.Fatalf("Length of getRes(%d) and num(%d) not equal!", len(getRes.Data), num)
	}

	for i := 0; i < num; i++ {
		dp := putReq.Data[i]
		if dp.Value.Value != getRes.Data[0].Values[i].Value ||
			dp.Value.Timestamp != getRes.Data[0].Values[i].Timestamp {
			t.Fatal("wrong result")
		}
	}
}

func dataGenerator(begin int64, numOfKeys, num int, shardId int32) *PutRequest {
	req := &PutRequest{}
	req.Data = make([]*tsdb.DataPoint, num*numOfKeys)
	index := 0

	for i := 0; i < numOfKeys; i++ {
		testKey := []byte(tsdb.RandStr(20))
		var testTime = begin

		for j := 0; j < num; j++ {
			req.Data[index] = &tsdb.DataPoint{
				Key: &tsdb.Key{
					Key:     testKey,
					ShardId: shardId,
				},
				Value: &tsdb.TimeValuePair{
					Timestamp: testTime,
					Value:     float64(100 + rand.Intn(50)),
				},
			}
			index++
			testTime += int64(55 + rand.Intn(10))
		}
	}
	return req
}

/*
func TestReload(t *testing.T) {
	tm := test_init()
	tm.start(s)
	defer tm.stop(s)

	s.Conf.Configer.Set(falcon.APP_CONF_FILE, map[string]string{
		"shardIds": "1,3,5,7",
	})
	time.Sleep(100 * time.Millisecond)
	tm.reload(s)
	time.Sleep(100 * time.Millisecond)

	var newMap []int
	for k, v := range tm.buckets {
		if v.GetState() == tsdb.OWNED {
			newMap = append(newMap, k)
		}
	}
	sort.Ints(newMap)
	wantMap := []int{1, 3, 5, 7}
	fmt.Println("newMap: ", newMap)
	fmt.Println("wantMap: ", wantMap)

	if len(wantMap) != len(newMap) {
		t.Fatal("wrong result of reload()")
	}
	for i := 0; i < len(wantMap); i++ {
		if wantMap[i] != newMap[i] {
			t.Fatal("wrong result of reload()")
		}
	}
}
*/
