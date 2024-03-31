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
		"CurrentConnections": graphdef.DefaultGraph("CurrentConnections", "float"),
		"ActiveConnections":  graphdef.DefaultGraph("ActiveConnections", "float"),
		"Statements":         graphdef.DefaultGraph("Statements", "float"),
		"StatementLatency":   graphdef.DefaultGraph("StatementLatency", "milliseconds"),
		"CPUUtilization":     graphdef.DefaultGraph("CPUUtilization", "percentage"),
		"MemoryUtilization":  graphdef.DefaultGraph("MemoryUtilization", "percentage"),
		"Traffic": {
			Label: "Traffic",
			Unit:  "bytes/sec",
			Metrics: []graphdef.Metric{
				{
					Name:  "NetworkReceiveBytes",
					Label: "Received Bytes",
				},
				{
					Name:  "NetworkTransmitBytes",
					Label: "Transmit Bytes",
				},
			},
		},
		"DbVolumeOperations": {
			Label: "Volume Operations",
			Unit:  "iops",
			Metrics: []graphdef.Metric{
				{
					Name:  "DbVolumeReadOperations",
					Label: "Read",
				},
				{
					Name:  "DbVolumeWriteOperations",
					Label: "Write",
				},
			},
		},
		"DbVolumeBytes": {
			Label: "Volume Bytes",
			Unit:  "bytes/sec",
			Metrics: []graphdef.Metric{
				{
					Name:  "DbVolumeReadBytes",
					Label: "Read",
				},
				{
					Name:  "DbVolumeWriteBytes",
					Label: "Write",
				},
			},
		},
		"DbVolumeUtilization": graphdef.DefaultGraph("DbVolumeUtilization", "percentage"),
	},
}

func main() {
	ctx := context.Background()

	optPrefix := flag.String("metric-key-prefix", "ocimds", "Metric key prefix")
	compartmentId := flag.String("compartmentId", "", "compartmentId")
	resourceId := flag.String("resourceId", "", "resourceId")
	optTitlePrefix := flag.String("title-prefix", "", "Title prefix")
	flag.Parse()

	if os.Getenv("MACKEREL_AGENT_PLUGIN_META") != "" {
		graphdef.Output(*optPrefix, *optTitlePrefix, def)
		os.Exit(0)
	}

	if *compartmentId == "" {
		log.Fatal("need -compartmentId")
	}
	if *resourceId == "" {
		log.Fatal("need -resourceId")
	}

	namespace := "oci_mysql_database"

	fillin := func(c string) string {
		return strings.Replace(c, "RESOURCE_ID", *resourceId, -1)
	}

	// https://docs.oracle.com/en-us/iaas/Content/Monitoring/Reference/mql.htm
	// https://docs.oracle.com/en-us/iaas/mysql-database/doc/metrics.html
	queries := []query.Query{
		{
			AggrigationName: "CurrentConnections",
			Query:           fillin(`CurrentConnections[1m]{resourceId="RESOURCE_ID"}.max()`),
		},
		{
			AggrigationName: "ActiveConnections",
			Query:           fillin(`ActiveConnections[1m]{resourceId="RESOURCE_ID"}.max()`),
		},
		{
			AggrigationName: "Statements",
			Query:           fillin(`Statements[1m]{resourceId="RESOURCE_ID"}.rate()`),
		},
		{
			AggrigationName: "StatementLatency",
			Query:           fillin(`StatementLatency[1m]{resourceId="RESOURCE_ID"}.rate()`),
		},
		{
			AggrigationName: "CPUUtilization",
			Query:           fillin(`CPUUtilization[1m]{resourceId="RESOURCE_ID"}.grouping().mean()`),
		},
		{
			AggrigationName: "MemoryUtilization",
			Query:           fillin(`MemoryUtilization[1m]{resourceId="RESOURCE_ID"}.grouping().mean()`),
		},
		{
			AggrigationName: "Traffic",
			Query:           fillin(`NetworkReceiveBytes[1m]{resourceId="RESOURCE_ID"}.rate().grouping().mean()`),
		},
		{
			AggrigationName: "Traffic",
			Query:           fillin(`NetworkTransmitBytes[1m]{resourceId="RESOURCE_ID"}.rate().grouping().mean()`),
		},
		{
			AggrigationName: "DbVolumeOperations",
			Query:           fillin(`DbVolumeReadOperations[1m]{resourceId="RESOURCE_ID"}.grouping().rate()`),
		},
		{
			AggrigationName: "DbVolumeOperations",
			Query:           fillin(`DbVolumeWriteOperations[1m]{resourceId="RESOURCE_ID"}.grouping().rate()`),
		},
		{
			AggrigationName: "DbVolumeBytes",
			Query:           fillin(`DbVolumeReadBytes[1m]{resourceId="RESOURCE_ID"}.grouping().rate()`),
		},
		{
			AggrigationName: "DbVolumeBytes",
			Query:           fillin(`DbVolumeWriteBytes[1m]{resourceId="RESOURCE_ID"}.grouping().rate()`),
		},
		{
			AggrigationName: "DbVolumeUtilization",
			Query:           fillin(`DbVolumeUtilization[1m]{resourceId="RESOURCE_ID"}.grouping().max()`),
		},
	}

	err := query.Run(ctx, *optPrefix, *compartmentId, namespace, queries)
	if err != nil {
		log.Fatalln(err)
	}
}
