/*
 * yubo@yubo.org
 * 2015-12-04
 */
package timer

import (
	"log"
	"testing"
	"time"
)

func ticker_cb(data interface{}) {
	i := data.(*int)
	*i += 1
	log.Println("ticker_cb", *i)
}

func timer_cb(data interface{}) {
	log.Println(data.(string))
}

func Test_timer(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)

	i := int(100)
	j := int(200)
	log.Println("start")
	NewTicker(time.Second, ticker_cb, &i)
	NewTimer(time.Second*2, timer_cb, "hello 2")
	NewTimer(time.Second, timer_cb, "hello 1")
	NewTimer(time.Second*3, timer_cb, "hello 3")
	time.Sleep(time.Second * 5)
	NewTimer(time.Second, timer_cb, "hello 6")
	NewTicker2(time.Second*1, time.Second*3, ticker_cb, &j)
	time.Sleep(time.Second * 10)
}
