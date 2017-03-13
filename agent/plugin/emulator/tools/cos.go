package main

import (
	"fmt"
	"math"
)

func main() {
	var (
		n        = 100
		interval = 60
		x        = (2 * math.Pi) / float64(n)
	)
	fmt.Printf("%d %d ", n, interval)
	for i := 0; i < n; i++ {
		fmt.Printf("%f ", math.Cos(x*float64(i)))
	}
}
