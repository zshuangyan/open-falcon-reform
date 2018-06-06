package model

import (
	"fmt"
)

type UserDefinedMetric struct {
	Name      string
	Command   string
	Step      int64
	MetricType string
}

type UserDefinedMetricsResponse struct {
	Metrics   []*UserDefinedMetric
}

func (this *UserDefinedMetric) String() string {
	return fmt.Sprintf( "<Metric: %v, Command: %v, Step: %v, Type: %v>",
		this.Name, this.Command, this.Step, this.MetricType)
}

func (this *UserDefinedMetricsResponse) String() string {
	return fmt.Sprintf(
		"Metrics:%v",
		this.Metrics,
	)
}
