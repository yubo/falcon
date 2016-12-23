package consistent

import (
	"fmt"
	"math"
	"os"
	"testing"
)

const (
	TEST_STR      = "test"
	width, height = 600, 320
	max           = 0xffffffff
)

func fwf(fd *os.File, format string, a ...interface{}) {
	fd.WriteString(fmt.Sprintf(format, a...))
}

func testSvg(replicas int) {
	fd, _ := os.OpenFile(fmt.Sprintf("%d.svg", replicas), os.O_CREATE|os.O_RDWR, 0644)
	defer fd.Close()
	c := New()
	c.NumberOfReplicas = replicas
	c.Add(TEST_STR)
	fwf(fd, "<svg xmlns='http://www.w3.org/2000/svg' style='stroke: grey; "+
		"fill: white; stroke-width: 0.7' width='%d' height='%d'>\n"+
		"<rect width='%d' height='%d' style='stroke-width: 0' />\n",
		width, height, width, height)
	wlen := len(c.sortedHashes)
	hmax := uint32(0)
	for last, i := uint32(0), 0; i < wlen; i++ {
		n := c.sortedHashes[i] - last
		last = c.sortedHashes[i]
		if n > hmax {
			hmax = n
		}
	}
	fwf(fd, "<polyline points='%d,%d", 0, height)
	for last, i := uint32(0), 0; i < wlen; i++ {
		x := (float32(width) / float32(wlen+1)) * float32(i+1)
		n := c.sortedHashes[i] - last
		last = c.sortedHashes[i]
		fwf(fd, " %f,%f",
			x, float32(height)-(float32(n)/float32(hmax))*float32(height))
	}
	fwf(fd, "'/></svg>")
}

func testVariance(replicas int) {
	var (
		sum, last uint32
		sum2, _x  float64
		i         int
	)

	c := New()
	c.NumberOfReplicas = replicas
	c.Add(TEST_STR)

	n := len(c.sortedHashes) + 1

	for last, i = 0, 0; i < n-1; i++ {
		x := c.sortedHashes[i] - last
		last = c.sortedHashes[i]
		sum += x
	}
	sum += max - last
	_x = float64(sum) / float64(n)

	for last, i = 0, 0; i < n-1; i++ {
		x := c.sortedHashes[i] - last
		last = c.sortedHashes[i]
		sum2 += math.Pow(float64(x)-_x, 2.0)
	}
	sum2 += math.Pow(float64(max-last)-_x, 2.0)
	fmt.Printf("replicas:%d variance:%f\n", replicas, sum2/float64(n-1))
}

func TestSvg5(t *testing.T)    { testSvg(5) }
func TestSvg50(t *testing.T)   { testSvg(50) }
func TestSvg500(t *testing.T)  { testSvg(500) }
func TestSvg5000(t *testing.T) { testSvg(5000) }

func TestVariance5(t *testing.T)      { testVariance(5) }
func TestVariance50(t *testing.T)     { testVariance(50) }
func TestVariance500(t *testing.T)    { testVariance(500) }
func TestVariance5000(t *testing.T)   { testVariance(5000) }
func TestVariance50000(t *testing.T)  { testVariance(50000) }
func TestVariance500000(t *testing.T) { testVariance(500000) }
