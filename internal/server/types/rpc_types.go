package types

import (
	"fmt"
	"strconv"

	pb "github.com/xChygyNx/metrical/internal/proto"
)

const (
	GAUGE   = "gauge"   // тип метрики gauge.
	COUNTER = "counter" // тип метрики counter.
)

type MetricGetter struct {
	ID    string
	MType string
}

func (mg *MetricGetter) SaveMetricInSyncMap(value string) (*pb.Metric, error) {
	var response pb.Metric

	switch mg.MType {
	case GAUGE:
		finalValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value of %s metric: %s", mg.MType, value)
		}
		response.Value = finalValue
	case COUNTER:
		finalDelta, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value of %s metric: %s", mg.MType, value)
		}
		response.Delta = finalDelta
	default:
		return nil, fmt.Errorf("unknown type %s of metric %s", mg.MType, mg.ID)
	}
	response.Id = mg.ID
	response.MType = mg.MType
	return &response, nil
}
