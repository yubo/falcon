package persistentKeyList

import (
	"github.com/huangaz/tsdb/lib/testUtil"
	"testing"
)

var (
	dataDirectory = testUtil.DataDirectory_Test
)

func TestWriteAndRead(t *testing.T) {
	var shardId int64 = 7
	testUtil.PathCreate(shardId)
	defer testUtil.FileDelete()

	keys := NewPersistentKeyList(shardId, dataDirectory)

	called := false
	_, err := ReadKeys(shardId, dataDirectory, func(KeyItem) bool {
		called = true
		return true
	})

	if err != nil {
		t.Fatal(err)
	}

	if called == true {
		t.Fatal("wrong")
	}

	data1 := KeyItem{5, "hi", 1}
	data2 := KeyItem{4, "test", 2}
	data3 := KeyItem{7, "bye", 3}

	keys.AppendKey(data1)
	keys.AppendKey(data2)
	keys.AppendKey(data3)
	err = keys.flush(true)
	if err != nil {
		t.Fatal(err)
	}

	var datas []KeyItem

	ReadKeys(shardId, dataDirectory, func(item KeyItem) bool {
		datas = append(datas, item)
		return true
	})

	if datas[0].Id != 5 || datas[0].Key != "hi" || datas[0].Category != 1 {
		t.Fatal("Wrong data!")
	}

	if datas[1].Id != 4 || datas[1].Key != "test" || datas[1].Category != 2 {
		t.Fatal("Wrong data!")
	}

	if datas[2].Id != 7 || datas[2].Key != "bye" || datas[2].Category != 3 {
		t.Fatal("Wrong data!")
	}

	// Rewrite two keys.
	i := 0
	err = keys.Compact(func() KeyItem {
		if i < 2 {
			item := KeyItem{1, "test2", 15}
			i++
			return item
		}
		return KeyItem{0, "", 0}
	})
	if err != nil {
		t.Fatal(err)
	}

	keys.AppendKey(KeyItem{8, "test3", 122})
	keys.flush(true)

	datas = datas[:0]
	ReadKeys(shardId, dataDirectory, func(item KeyItem) bool {
		datas = append(datas, item)
		return true
	})

	if datas[0].Id != 1 || datas[0].Key != "test2" || datas[0].Category != 15 {
		t.Fatal("wrong data!")
	}

	if datas[2].Id != 8 || datas[2].Key != "test3" || datas[2].Category != 122 {
		t.Fatal("wrong data!")
	}

}
