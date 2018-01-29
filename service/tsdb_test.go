package service

import (
	"math/rand"
	"testing"
	"time"

	"github.com/huangaz/tsdb/lib/testUtil"
	"github.com/yubo/falcon"
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

	tsdb := &TsdbModule{}
	tsdb.prestart(s)

	return tsdb
}

/*
func TestPut(t *testing.T) {
	test_init()
	tsdb.start(s)
	defer tsdb.stop(s)

	num := 1000
	putReq := dataGenerator(1, num)

	time.Sleep(100 * time.Millisecond)

	putRes, err := tsdb.put(putReq)
	if err != nil {
		t.Fatal(err)
	}
	if putRes.N != int32(num) {
		t.Fatalf("Put failed! Want %d,put %d", num, putRes.N)
	}

	time.Sleep(100 * time.Millisecond)

	getReq := &GetRequest{
		Start:   0,
		End:     int64(60 * num),
		ShardId: 1,
		Key:     putReq.Items[0].Key,
	}
	getRes, err := tsdb.get(getReq)
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Println(putReq.PrintForDebug())
	// fmt.Println(getRes.PrintForDebug())

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
*/

func TestPutAndRecover(t *testing.T) {
	tsdb := test_init()
	tsdb.start(s)

	num := 1000
	putReq := dataGenerator(time.Now().Unix()-int64(num*60), 1, num)

	time.Sleep(100 * time.Millisecond)

	putRes, err := tsdb.put(putReq)
	if err != nil {
		t.Fatal(err)
	}
	if putRes.N != int32(num) {
		t.Fatalf("Put failed! Want %d,put %d", num, putRes.N)
	}

	time.Sleep(1000 * time.Millisecond)

	// finalizedTimstamp := time.Now().Unix() - int64(allowedTimestampBehind+60) - int64(bucketUtils.Duration(1, bucketSize))
	n, err := tsdb.buckets[1].FinalizeBuckets(0)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatal("finalize failed")
	}
	// log.Println("lastFinalizedBucket_: ", tsdb.buckets[1].GetLastFinalizedBucket())

	tsdb.stop(s)

	// recover data from disk
	t2 := test_init()
	t2.start(s)
	time.Sleep(100 * time.Millisecond)

	getReq := &GetRequest{
		Start: 0,
		End:   time.Now().Unix(),
		Keys: []*Key{
			&Key{
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

	if len(putReq.Data) != len(getRes.Data) {
		t.Fatalf("Length of putReq(%d) and getRes(%d) not equal!", len(putReq.Data), len(getRes.Data))
	}

	for i, dp := range putReq.Data {
		if dp.Value.Value != getRes.Data[0].Values[i].Value ||
			dp.Value.Timestamp != getRes.Data[0].Values[i].Timestamp {
			t.Fatal("wrong result")
		}
	}
}

func dataGenerator(begin int64, numOfKeys, num int) *PutRequest {
	req := &PutRequest{}
	req.Data = make([]*DataPoint, num*numOfKeys)
	index := 0

	for i := 0; i < numOfKeys; i++ {
		testKey := []byte(testUtil.RandStr(10))
		var testTime = begin

		for j := 0; j < num; j++ {
			req.Data[index] = &DataPoint{
				Key: &Key{
					Key:     testKey,
					ShardId: int32(i + 1),
				},
				Value: &TimeValuePair{
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
	test_init()
	tsdb.start(s)
	defer tsdb.stop(s)

	s.Conf.Configer.Set(APP_CONF_FILE, map[string]string{
		"shardIds": "1,3,5,7",
	})
	time.Sleep(100 * time.Millisecond)
	tsdb.reload(s)
	time.Sleep(100 * time.Millisecond)

	var newMap []int
	for k, v := range tsdb.buckets {
		if v.GetState() == bucketMap.OWNED {
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
