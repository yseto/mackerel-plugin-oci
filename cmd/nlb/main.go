package main

import (
	"cmp"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yseto/mackerel-plugin-oci/internal/graphdef"
	"github.com/yseto/mackerel-plugin-oci/internal/query"
)

var def = graphdef.Def{
	Graphs: map[string]graphdef.Graph{
		"ProcessedBytes":   graphdef.DefaultGraph("ProcessedBytes", "bytes"),
		"ProcessedPackets": graphdef.DefaultGraph("ProcessedPackets", "float"),
		"DroppedBySecurityLists": {
			Label: "Dropped by Security lists",
			Unit:  "float",
			Metrics: []graphdef.Metric{
				{
					Name:  "IngressPacketsDroppedBySL",
					Label: "Ingress Packets",
				},
				{
					Name:  "EgressPacketsDroppedBySL",
					Label: "Egress Packets",
				},
			},
		},
		"Backends": {
			Label: "Backends",
			Unit:  "float",
			Metrics: []graphdef.Metric{
				{
					Name:  "HealthyBackendsPerNlb",
					Label: "Healthy",
				},
				{
					Name:  "UnhealthyBackendsPerNlb",
					Label: "Unhealthy",
				},
			},
		},
		"NewConnections": {
			Label: "NewConnections",
			Unit:  "float",
			Metrics: []graphdef.Metric{
				{
					Name:  "NewConnections",
					Label: "New",
				},
				{
					Name:  "NewConnectionsTCP",
					Label: "TCP",
				},
				{
					Name:  "NewConnectionsUDP",
					Label: "UDP",
				},
			},
		},
	},
}

func main() {
	ctx := context.Background()

	optPrefix := flag.String("metric-key-prefix", "ocinlb", "Metric key prefix")
	compartmentId := flag.String("compartmentId", "", "compartmentId")
	resourceId := flag.String("resourceId", "", "resourceId")
	resourceName := flag.String("resourceName", "", "resourceName")
	optTitlePrefix := flag.String("title-prefix", "", "Title prefix")
	flag.Parse()

	if *compartmentId == "" {
		log.Fatal("need -compartmentId")
	}
	if *resourceId == "" {
		log.Fatal("need -resourceId")
	}
	if *resourceName == "" {
		log.Fatal("need -resourceName")
	}

	prefix := fmt.Sprintf("%s.%s", *optPrefix, *resourceName)

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		graphdef.Output(prefix, cmp.Or(*optTitlePrefix, *resourceName), def)
		os.Exit(0)
	}

	namespace := "oci_nlb"

	fillin := func(c string) string {
		c = strings.Replace(c, "RESOURCE_NAME", *resourceName, -1)
		return strings.Replace(c, "RESOURCE_ID", *resourceId, -1)
	}

	// https://docs.oracle.com/en-us/iaas/Content/Monitoring/Reference/mql.htm
	// https://docs.oracle.com/en-us/iaas/Content/NetworkLoadBalancer/Metrics/metrics.htm
	queries := []query.Query{
		{
			AggrigationName: "ProcessedBytes",
			Query:           fillin(`ProcessedBytes[1m]{resourceId = "RESOURCE_ID", resourceName = "RESOURCE_NAME"}.sum()`),
		},
		{
			AggrigationName: "ProcessedPackets",
			Query:           fillin(`ProcessedPackets[1m]{resourceId = "RESOURCE_ID", resourceName = "RESOURCE_NAME"}.sum()`),
		},
		{
			AggrigationName: "DroppedBySecurityLists",
			Query:           fillin(`IngressPacketsDroppedBySL[1m]{resourceId = "RESOURCE_ID", resourceName = "RESOURCE_NAME"}.sum()`),
		},
		{
			AggrigationName: "DroppedBySecurityLists",
			Query:           fillin(`EgressPacketsDroppedBySL[1m]{resourceId = "RESOURCE_ID", resourceName = "RESOURCE_NAME"}.sum()`),
		},
		{
			AggrigationName: "Backends",
			Query:           fillin(`HealthyBackendsPerNlb[1m]{resourceId = "RESOURCE_ID", resourceName = "RESOURCE_NAME"}.max()`),
		},
		{
			AggrigationName: "Backends",
			Query:           fillin(`UnhealthyBackendsPerNlb[1m]{resourceId = "RESOURCE_ID", resourceName = "RESOURCE_NAME"}.max()`),
		},
		{
			AggrigationName: "NewConnections",
			Query:           fillin(`NewConnections[1m]{resourceId = "RESOURCE_ID", resourceName = "RESOURCE_NAME"}.sum()`),
		},
		{
			AggrigationName: "NewConnections",
			Query:           fillin(`NewConnectionsTCP[1m]{resourceId = "RESOURCE_ID", resourceName = "RESOURCE_NAME"}.sum()`),
		},
		{
			AggrigationName: "NewConnections",
			Query:           fillin(`NewConnectionsUDP[1m]{resourceId = "RESOURCE_ID", resourceName = "RESOURCE_NAME"}.sum()`),
		},
	}

	err := query.Run(ctx, prefix, *compartmentId, namespace, queries)
	if err != nil {
		log.Fatalln(err)
	}
}
