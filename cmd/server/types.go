package server

type gauge float64

type counter int64

type PollCount counter

type RandomValue gauge

type MemStorage struct {
	gauges map[string]gauge			`json:"gauges"`
	counters map[string]counter		`json:"counters"`
}

var instance *MemStorage

func GetMemStorage() *MemStorage {
	if instance == nil {
		instance = &MemStorage{}
	}
	return instance
}

func (ms *MemStorage) SetGauges(data map[string]float64) {
	for k, v := range data {
		ms.gauges[k] = gauge(v)
	}
}

func (ms *MemStorage) SetConunters(data map[string]float64) {
	for k, v := range data {
		ms.counters[k] += counter(v)
	}
}

