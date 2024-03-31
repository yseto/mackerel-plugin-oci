package query

import (
	"cmp"
	"context"
	"fmt"

	"github.com/yseto/mackerel-plugin-oci/internal/mql"
)

type Query struct {
	AggrigationName          string
	AggrigationDimensionName string

	Query         string
	DimensionName string

	Scale float64
}

func Run(ctx context.Context, prefix, compartmentId, namespace string, queries []Query) error {
	handler, err := mql.NewHandler()
	if err != nil {
		return err
	}

	for _, q := range queries {
		result, err := handler.QueryWithBackoffRetry(ctx, mql.QueryInput{
			CompartmentId: &compartmentId,
			Namespace:     &namespace,
			Query:         &q.Query,
		})
		if err != nil {
			return err
		}

		for _, item := range result.Items {
			dp := item.Datapoints[len(item.Datapoints)-1]

			name := item.Name
			if q.DimensionName != "" {
				name = cmp.Or(item.Dimensions[q.DimensionName], "undefined")
			}

			aggrigationName := q.AggrigationName
			if q.AggrigationDimensionName != "" {
				aggrigationName = fmt.Sprintf("%s.%s",
					q.AggrigationName,
					cmp.Or(item.Dimensions[q.AggrigationDimensionName], "undefined"),
				)
			}

			fmt.Printf("%s.%s.%s\t%f\t%d\n", prefix, aggrigationName, name, (dp.Value * cmp.Or(q.Scale, 1.0)), dp.Time.Unix())
		}
	}

	return nil
}
