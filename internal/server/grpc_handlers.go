package server

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/xChygyNx/metrical/internal/proto"
	"github.com/xChygyNx/metrical/internal/server/types"
)

type MetricServer struct {
	proto.BatchMetricHandlerServer

	metrics sync.Map
}

func (ms *MetricServer) saveGauge(id string, value float64) {
	valueStr := strconv.FormatFloat(value, 'f', -1, 64)
	ms.metrics.Store(id, valueStr)
}

func (ms *MetricServer) saveCounter(id string, value int64) (newValue int64, err error) {
	var saveValue string
	prevValue, ok := ms.metrics.LoadAndDelete(id)
	if !ok {
		saveValue = strconv.FormatInt(value, 64)
	} else {
		prevValueStr, ok := prevValue.(string)
		if !ok {
			return 0, fmt.Errorf("invalid value in metric store: %v", prevValue)
		}
		prevValueInt, err := strconv.ParseInt(prevValueStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid value in metric store: %v", prevValueStr)
		}
		newValue = prevValueInt + value
		saveValue = strconv.FormatInt(newValue, 64)
	}
	ms.metrics.Store(id, saveValue)
	return newValue, nil
}

func (ms *MetricServer) SaveBatchMetrics(
	ctx context.Context, in *proto.BatchMetricRequest) (*proto.BatchMetricResponse, error) {
	var response proto.BatchMetricResponse

	for ind, metric := range in.Metrics {
		switch metric.GetMType() {

		case GAUGE:
			ms.saveGauge(metric.GetId(), metric.GetValue())
		case COUNTER:
			newMetricValue, err := ms.saveCounter(metric.GetId(), metric.GetDelta())
			if err != nil {
				return nil, fmt.Errorf("error in save counter metric: %w", err)
			}
			in.Metrics[ind].Delta = newMetricValue
		default:
			return nil, fmt.Errorf("invalid type of metric %s: %s", metric.Id, metric.MType)
		}
	}

	response.Status = proto.Status_OK
	response.Metrics = in.Metrics
	return &response, nil
}

func (ms *MetricServer) GetJSONMetric(
	ctx context.Context, in *proto.GetJSONMetricsRequest) (*proto.GetJSONMetricResponse, error) {
	var result proto.GetJSONMetricResponse

	responseGetter := types.MetricGetter{
		ID:    in.GetId(),
		MType: in.GetMType(),
	}
	value, ok := ms.metrics.Load(in.GetId())
	if !ok {
		return nil, fmt.Errorf("not found metric %s", in.GetId())
	}
	valueStr, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("invalid value of %s: %v", in.GetId(), value)
	}

	response, err := responseGetter.SaveMetricInSyncMap(valueStr)
	if err != nil {
		result.Status = proto.Status_ERROR
		return nil, fmt.Errorf("error in getting value of metric %s: %w", in.GetId(), err)
	} else {
		result.Status = proto.Status_OK
		result.Metric = response
	}

	return &result, nil

}
