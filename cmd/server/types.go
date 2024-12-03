package main

type gauge float64

type counter int64

type memStorage struct {
	Gauges   map[string]gauge   `json:"gauges"`
	Counters map[string]counter `json:"counters"`
}

var instance *memStorage

func GetMemStorage() *memStorage {
	if instance == nil {
		instance = &memStorage{}
		instance.Gauges = make(map[string]gauge)
		instance.Counters = make(map[string]counter)
	}
	return instance
}

func (ms *memStorage) SetGauge(mName string, mValue float64) {
	ms.Gauges[mName] = gauge(mValue)
}

func (ms *memStorage) SetConunter(mName string, mValue int64) {
	ms.Counters[mName] += counter(mValue)
}

func (ms *memStorage) SetGauges(data map[string]float64) {
	for k, v := range data {
		ms.Gauges[k] = gauge(v)
	}
}

func (ms *memStorage) GetGauge(mName string) (float64, bool) {
	metric, ok := ms.Gauges[mName]
	return float64(metric), ok
}

func (ms *memStorage) GetCounter(mName string) (int64, bool) {
	metric, ok := ms.Counters[mName]
	return int64(metric), ok
}

func (ms *memStorage) SetConunters(data map[string]float64) {
	for k, v := range data {
		ms.Counters[k] += counter(v)
	}
}
