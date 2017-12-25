package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	var (
		n        = 100
		interval = 60
		x        = (2 * math.Pi) / float64(n)
		offset   = 0
	)
	if len(os.Args) == 2 {
		if i, err := strconv.Atoi(os.Args[1]); err == nil {
			offset = i
		}
	}
	fmt.Printf("%d %d ", n, interval)
	for i := offset; i < n+offset; i++ {
		fmt.Printf("%f ", math.Cos(x*float64(i)))
	}
}
