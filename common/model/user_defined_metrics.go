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

type UserDefinedMetricHost struct{
	Name       string   `json:"name" gorm:"column:name"`
	Command    string   `json:"command" gorm:"column:command"`
	Step       int      `json:"step" gorm:"column:step"`
	MetricType string   `json:"metric_type" gorm:"column:metric_type"`
	ValueType  string   `json:"value_type" gorm:"column:value_type"`
	HostID     int64    `json:"host_id" gorm:"column:host_id"`
}

func (this UserDefinedMetricHost) TableName() string {
	return "user_defined_metric"
}