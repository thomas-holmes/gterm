package main

import (
	"log"
	"math"
	"time"
)

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
func max64(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
func maxu64(a uint64, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
func maxf32(a float32, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
func min64(a int64, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
func minu64(a uint64, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
func minf32(a float32, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -1 * a
}

func timeMe(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%v %s", name, elapsed)

}

func distance(x0 int, y0 int, x1 int, y1 int) float64 {
	x := x1 - x0
	y := y1 - y0
	return math.Sqrt(float64(x*x) + float64(y*y))
}
