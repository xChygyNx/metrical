package types

import (
	"fmt"
	"strconv"
	"sync"

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

func (mg *MetricGetter) GetJSONResponse(storage sync.Map) (*pb.Metric, error) {
	var response *pb.Metric

	value, ok := storage.Load(mg.ID)
	if !ok {
		return nil, fmt.Errorf("not found metric %s", mg.ID)
	}
	valueStr, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("invalid value of %s: %v", mg.ID, value)
	}

	switch mg.MType {
	case GAUGE:
		finalValue, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value of %s metric: %s", mg.MType, valueStr)
		}
		response.Value = finalValue
	case COUNTER:
		finalDelta, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value of %s metric: %s", mg.MType, valueStr)
		}
		response.Delta = finalDelta
	default:
		return nil, fmt.Errorf("unknown type %s of metric %s", mg.MType, mg.ID)
	}
	response.Id = mg.ID
	response.MType = mg.MType
	return response, nil
}
