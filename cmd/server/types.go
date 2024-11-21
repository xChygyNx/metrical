package main

type gauge float64

type counter int64

type MemStorage struct {
	gauges map[string]gauge
	counters map[string]counter
}