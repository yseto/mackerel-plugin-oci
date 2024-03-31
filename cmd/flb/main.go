package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/yseto/mackerel-plugin-oci/internal/graphdef"
	"github.com/yseto/mackerel-plugin-oci/internal/query"
)

var def = graphdef.Def{
	Graphs: map[string]graphdef.Graph{
		"HttpRequests":         graphdef.DefaultGraph("HttpRequests", "float"),
		"ActiveConnections":    graphdef.DefaultGraph("ActiveConnections", "float"),
		"ActiveSSLConnections": graphdef.DefaultGraph("ActiveSSLConnections", "float"),
		"Traffic.#": {
			Label: "Traffic",
			Unit:  "bytes/sec",
			Metrics: []graphdef.Metric{
				{
					Name:  "BytesReceived",
					Label: "Received Bytes",
				},
				{
					Name:  "BytesSent",
					Label: "Sent Bytes",
				},
			},
		},
		"AcceptedConnections": graphdef.DefaultGraph("AcceptedConnections", "float"),
		"HandledConnections":  graphdef.DefaultGraph("HandledConnections", "float"),
		"SSLHandshake": {
			Label: "SSLHandshake",
			Unit:  "float",
			Metrics: []graphdef.Metric{
				{
					Name:  "FailedSSLHandshake",
					Label: "Failed",
				},
				{
					Name:  "AcceptedSSLHandshake",
					Label: "Accepted",
				},
				{
					Name:  "FailedSSLClientCertVerify",
					Label: "FailedSSLClientCertVerify",
				},
			},
		},
		"PeakBandwidth": graphdef.DefaultGraph("PeakBandwidth", "bits/sec"),
		"HttpResponses.#": {
			Label: "HttpResponses",
			Unit:  "float",
			Metrics: []graphdef.Metric{
				{
					Name:  "HttpResponses2xx",
					Label: "2xx",
				},
				{
					Name:  "HttpResponses200",
					Label: "200",
				},
				{
					Name:  "HttpResponses3xx",
					Label: "3xx",
				},
				{
					Name:  "HttpResponses4xx",
					Label: "4xx",
				},
				{
					Name:  "HttpResponses502",
					Label: "502",
				},
				{
					Name:  "HttpResponses504",
					Label: "504",
				},
				{
					Name:  "HttpResponses5xx",
					Label: "5xx",
				},
			},
		},
	},
}

func main() {
	ctx := context.Background()

	optPrefix := flag.String("metric-key-prefix", "ociflb", "Metric key prefix")
	compartmentId := flag.String("compartmentId", "", "compartmentId")
	resourceId := flag.String("resourceId", "", "resourceId")
	optTitlePrefix := flag.String("title-prefix", "", "Title prefix")
	flag.Parse()

	if *compartmentId == "" {
		log.Fatal("need -compartmentId")
	}
	if *resourceId == "" {
		log.Fatal("need -resourceId")
	}

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		graphdef.Output(*optPrefix, *optTitlePrefix, def)
		os.Exit(0)
	}

	namespace := "oci_lbaas"

	fillin := func(c string) string {
		return strings.Replace(c, "RESOURCE_ID", *resourceId, -1)
	}

	// https://docs.oracle.com/en-us/iaas/Content/Monitoring/Reference/mql.htm
	// https://docs.oracle.com/en-us/iaas/Content/Balance/Reference/loadbalancermetrics.htm
	queries := []query.Query{
		{
			AggrigationName: "HttpRequests",
			Query:           fillin(`HttpRequests[1m]{resourceId = "RESOURCE_ID"}.grouping().sum()`),
		},
		{
			AggrigationName: "ActiveConnections",
			Query:           fillin(`ActiveConnections[1m]{resourceId = "RESOURCE_ID"}.grouping().sum()`),
		},
		{
			AggrigationName: "ActiveSSLConnections",
			Query:           fillin(`ActiveSSLConnections[1m]{resourceId = "RESOURCE_ID"}.grouping().sum()`),
		},
		{
			AggrigationName:          "Traffic",
			AggrigationDimensionName: "lbComponent",
			Query:                    fillin(`BytesReceived[1m]{resourceId = "RESOURCE_ID"}.groupBy(lbComponent).sum()`),
			Scale:                    (1.0 / 60),
		},
		{
			AggrigationName:          "Traffic",
			AggrigationDimensionName: "lbComponent",
			Query:                    fillin(`BytesSent[1m]{resourceId = "RESOURCE_ID" }.groupBy(lbComponent).sum()`),
			Scale:                    (1.0 / 60),
		},
		{
			AggrigationName: "AcceptedConnections",
			Query:           fillin(`AcceptedConnections[1m]{resourceId = "RESOURCE_ID"}.grouping().sum()`),
		},
		{
			AggrigationName: "HandledConnections",
			Query:           fillin(`HandledConnections[1m]{resourceId = "RESOURCE_ID"}.grouping().sum()`),
		},
		{
			AggrigationName: "SSLHandshake",
			Query:           fillin(`FailedSSLHandshake[1m]{resourceId = "RESOURCE_ID"}.grouping().sum()`),
		},
		{
			AggrigationName: "SSLHandshake",
			Query:           fillin(`AcceptedSSLHandshake[1m]{resourceId = "RESOURCE_ID"}.grouping().sum()`),
		},
		{
			AggrigationName: "SSLHandshake",
			Query:           fillin(`FailedSSLClientCertVerify[1m]{resourceId = "RESOURCE_ID"}.grouping().sum()`),
		},
		{
			AggrigationName: "PeakBandwidth",
			Query:           fillin(`PeakBandwidth[1m]{resourceId = "RESOURCE_ID"}.grouping().sum()`),
			Scale:           1024 * 1024,
		},
		{
			AggrigationName:          "HttpResponses",
			AggrigationDimensionName: "listenerName",
			Query:                    fillin(`HttpResponses200[1m]{resourceId = "RESOURCE_ID", lbComponent = "Listener"}.groupBy(listenerName).sum()`),
		},
		{
			AggrigationName:          "HttpResponses",
			AggrigationDimensionName: "listenerName",
			Query:                    fillin(`HttpResponses2xx[1m]{resourceId = "RESOURCE_ID", lbComponent = "Listener"}.groupBy(listenerName).sum()`),
		},
		{
			AggrigationName:          "HttpResponses",
			AggrigationDimensionName: "listenerName",
			Query:                    fillin(`HttpResponses3xx[1m]{resourceId = "RESOURCE_ID", lbComponent = "Listener"}.groupBy(listenerName).sum()`),
		},
		{
			AggrigationName:          "HttpResponses",
			AggrigationDimensionName: "listenerName",
			Query:                    fillin(`HttpResponses4xx[1m]{resourceId = "RESOURCE_ID", lbComponent = "Listener"}.groupBy(listenerName).sum()`),
		},
		{
			AggrigationName:          "HttpResponses",
			AggrigationDimensionName: "listenerName",
			Query:                    fillin(`HttpResponses502[1m]{resourceId = "RESOURCE_ID", lbComponent = "Listener"}.groupBy(listenerName).sum()`),
		},
		{
			AggrigationName:          "HttpResponses",
			AggrigationDimensionName: "listenerName",
			Query:                    fillin(`HttpResponses504[1m]{resourceId = "RESOURCE_ID", lbComponent = "Listener"}.groupBy(listenerName).sum()`),
		},
		{
			AggrigationName:          "HttpResponses",
			AggrigationDimensionName: "listenerName",
			Query:                    fillin(`HttpResponses5xx[1m]{resourceId = "RESOURCE_ID", lbComponent = "Listener"}.groupBy(listenerName).sum()`),
		},
	}

	err := query.Run(ctx, *optPrefix, *compartmentId, namespace, queries)
	if err != nil {
		log.Fatalln(err)
	}
}
