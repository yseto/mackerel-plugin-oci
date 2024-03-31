package mql

import (
	"context"
	"time"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/monitoring"
)

type MonitoringClient interface {
	SummarizeMetricsData(ctx context.Context, request monitoring.SummarizeMetricsDataRequest) (response monitoring.SummarizeMetricsDataResponse, err error)
}

type Handler struct {
	client             MonitoringClient
	StartTime, EndTime time.Time
}

func NewHandler() (*Handler, error) {
	provider := common.DefaultConfigProvider()
	client, err := monitoring.NewMonitoringClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, err
	}

	endTime := time.Now()
	startTime := endTime.Add(-2 * time.Minute)

	return &Handler{client: &client, StartTime: startTime, EndTime: endTime}, nil
}

type QueryInput struct {
	CompartmentId, Namespace, Query, ResourceGroup, Resolution *string
}

type DataPoint struct {
	Time  time.Time
	Value float64
}

type QueryResultItem struct {
	Name       string
	Datapoints []*DataPoint
	Dimensions map[string]string
}

type QueryResult struct {
	Items []*QueryResultItem
}

func (h *Handler) Query(ctx context.Context, input QueryInput) (*QueryResult, error) {
	req := monitoring.SummarizeMetricsDataRequest{
		CompartmentId: input.CompartmentId,
		SummarizeMetricsDataDetails: monitoring.SummarizeMetricsDataDetails{
			Namespace:     input.Namespace,
			Query:         input.Query,
			ResourceGroup: input.ResourceGroup,
			StartTime:     &common.SDKTime{Time: h.StartTime},
			EndTime:       &common.SDKTime{Time: h.EndTime},
			Resolution:    input.Resolution,
		},
	}
	resp, err := h.client.SummarizeMetricsData(ctx, req)
	if err != nil {
		return nil, err
	}

	result := &QueryResult{}
	for _, v := range resp.Items {
		if v.Name == nil {
			continue
		}
		item := &QueryResultItem{
			Name:       *v.Name,
			Dimensions: v.Dimensions,
		}
		for _, vv := range v.AggregatedDatapoints {
			if vv.Value == nil {
				continue
			}
			dp := &DataPoint{
				Time:  vv.Timestamp.Time,
				Value: *vv.Value,
			}
			item.Datapoints = append(item.Datapoints, dp)
		}
		if len(item.Datapoints) > 0 {
			result.Items = append(result.Items, item)
		}
	}
	return result, nil
}
