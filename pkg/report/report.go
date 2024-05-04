package report

import (
	"sort"
	"time"
)

type Result struct {
	TotalTime time.Duration
	Err       error
}

type Percentile struct {
	Percentile float64
	Latency    time.Duration
}

type Stats struct {
	TotalTime   time.Duration
	Total       int
	Fastest     time.Duration
	Slowest     time.Duration
	Average     time.Duration
	RPS         float64
	Percentiles []Percentile
	Errors      map[string]int
}

type Report interface {
	Results() chan<- Result
	Run() <-chan Stats
}

type report struct {
	results chan Result
	stats   Stats
}

func NewReport(totalClients uint) Report {
	return &report{
		results: make(chan Result, totalClients),
		stats: Stats{
			Errors: make(map[string]int),
		},
	}
}

func (r *report) Results() chan<- Result {
	return r.results
}

func (r *report) Run() <-chan Stats {
	donec := make(chan Stats, 1)
	go func() {
		defer close(donec)
		r.processResults()
		donec <- r.stats
	}()
	return donec
}

func (r *report) processResults() {
	start := time.Now()
	latencies := []time.Duration{}
	for res := range r.results {
		if res.Err != nil {
			r.stats.Errors[res.Err.Error()]++
			continue
		}
		latencies = append(latencies, res.TotalTime)
	}
	r.stats.TotalTime = time.Since(start)

	if len(latencies) == 0 {
		return
	}
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	r.stats.Total = len(latencies)
	r.stats.Fastest = latencies[0]
	r.stats.Slowest = latencies[len(latencies)-1]
	
	var avgTotal time.Duration
	for _, total := range latencies {
		avgTotal += total
	}
	r.stats.Average = time.Duration(int(avgTotal.Nanoseconds()) / len(latencies))

	r.stats.RPS = float64(len(latencies)) / r.stats.TotalTime.Seconds()

	for _, percentile := range []float64{10, 25, 50, 75, 90, 95, 99, 99.9} {
		i := int(float64(len(latencies)) * percentile / 100.0)
		r.stats.Percentiles = append(r.stats.Percentiles, Percentile{Percentile: percentile, Latency: latencies[i]})
	}
}
