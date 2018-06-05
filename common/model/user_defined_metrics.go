package model

import "fmt"

type UserDefinedMetric struct {
	Metric_Name  string
	Command      string
}

type UserDefinedMetricsResponse struct {
	Metrics   []*UserDefinedMetric
}

func (this *UserDefinedMetric) String() string {
	return fmt.Sprintf(
		"%s/%s",
		this.Metric_Name,
		this.Command,
	)
}

func (this *UserDefinedMetricsResponse) String() string {
	return fmt.Sprintf(
		"<Metrics:%v>",
		this.Metrics,
	)
}
