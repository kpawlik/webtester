package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type stat struct {
	mx               sync.Mutex
	delay            time.Duration
	times            []time.Duration
	successNo, total int64
}

func newStat(total int64, delay time.Duration) *stat {
	return &stat{
		times:     make([]time.Duration, 0),
		successNo: 0,
		total:     total,
		delay:     delay,
	}
}
func (s *stat) AddTime(duration time.Duration) {
	s.mx.Lock()
	s.times = append(s.times, duration)
	s.mx.Unlock()
}
func (s *stat) Success() {
	s.mx.Lock()
	s.successNo++
	s.mx.Unlock()
}
func (s *stat) AvgTime() float64 {
	var sum time.Duration
	for _, duration := range s.times {
		sum += duration
	}
	return sum.Seconds() / float64(s.successNo)
}
func (s *stat) CalcApprox() (approx time.Duration) {
	seconds := math.Ceil(float64(s.total) / (float64(time.Second) / float64(s.delay.Nanoseconds())))
	approx, _ = time.ParseDuration(fmt.Sprintf("%fs", seconds))
	return

}
func (s *stat) CalcRps() (rp int, unit string) {
	rps := float64(time.Second) / float64(wait.Nanoseconds())
	if rps > 1 {
		rp = int(math.Ceil(rps))
		unit = "second"
		return
	}
	rpm := math.Ceil(float64(time.Minute) / float64(wait.Nanoseconds()))
	rp = int(math.Ceil(rpm))
	unit = "minute"
	return
}

func (s *stat) MaxTime() time.Duration {
	var max time.Duration
	for _, t := range s.times {
		if t > max {
			max = t
		}
	}
	return max
}

func (s *stat) MinTime() time.Duration {
	var min time.Duration
	if len(s.times) > 0 {
		min = s.times[0]
	}
	for _, t := range s.times {
		if t < min {
			min = t
		}
	}
	return min
}
