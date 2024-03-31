package graphdef

import (
	"encoding/json"
	"fmt"
)

// refer:	github.com/mackerelio/go-mackerel-plugin v0.1.4

type Def struct {
	Graphs map[string]Graph `json:"graphs"`
}

type Graph struct {
	Label   string   `json:"label"`
	Unit    string   `json:"unit"`
	Metrics []Metric `json:"metrics"`
}

type Metric struct {
	Name    string `json:"name"`
	Label   string `json:"label"`
	Stacked bool   `json:"stacked"`
}

func Output(prefix, titlePrefix string, def Def) {
	defWithPrefix := Def{Graphs: map[string]Graph{}}

	for k, v := range def.Graphs {
		if titlePrefix != "" {
			v.Label = fmt.Sprintf("%s %s", titlePrefix, v.Label)
		}

		defWithPrefix.Graphs[fmt.Sprintf("%s.%s", prefix, k)] = v
	}

	b, _ := json.Marshal(defWithPrefix)
	fmt.Println("# mackerel-agent-plugin")
	fmt.Println(string(b))
}

func DefaultGraph(label, unit string) Graph {
	return Graph{
		Label: label,
		Unit:  unit,
		Metrics: []Metric{
			{
				Name:  "*",
				Label: "%1",
			},
		},
	}
}
