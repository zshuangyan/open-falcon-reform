package model

import (
	"fmt"
)

type AddedMetric struct {
	Name      string
	Command   string
	Step      int64
	MetricType string
}

type RemovedMetric struct {
	Name      string
}

type AddedMetricsResponse struct {
	Metrics   []*AddedMetric
}

type RemovedMetricsResponse struct {
	Metrics   []*RemovedMetric
}

func (this *AddedMetric) String() string {
	return fmt.Sprintf( "<Metric: %v, Command: %v, Step: %v, Type: %v>",
		this.Name, this.Command, this.Step, this.MetricType)
}

func (this *RemovedMetric) String() string {
	return fmt.Sprintf("<Metric: %v>", this.Name)
}

func (this *AddedMetricsResponse) String() string {
	return fmt.Sprintf(
		"Metrics:%v",
		this.Metrics,
	)
}

func (this *RemovedMetricsResponse) String() string {
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
	Status     int      `json:"status" gorm:"column:status"`
}

func (this UserDefinedMetricHost) TableName() string {
	return "user_defined_metric"
}