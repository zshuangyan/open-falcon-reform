// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cron

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"log"
	"time"
	"github.com/rakyll/ticktock"
	"github.com/rakyll/ticktock/t"
	"os/exec"
	"strconv"
	"strings"
)

type UserDefinedMetric struct {
	Name      string
	Command   string
	Step      int64
	MetricType string
}

func ExecCommand (c string) (int, error) {
	result, err := exec.Command("bash", "-c", c).Output()
	if err != nil {
		log.Println("Command finished with error:", err)
		return -1, err
	}
	data, err := strconv.Atoi(strings.TrimSpace(string(result[:])))
	if err != nil {
		log.Println("Command finished with error:", err)
		return -1, err
	}
	return data, nil
}

func (m *UserDefinedMetric) Run() error{
	mvs := []*model.MetricValue{}
	mv := new(model.MetricValue)
	mv.Metric = m.Name
	mv.Step = m.Step * 60
	mv.Type = m.MetricType
	mv.Timestamp = time.Now().Unix()
	log.Printf("\nUserDefinedMetric: %v\nCommand: %v\nStep: %v\n", m.Name, m.Command, m.Step)
	value, err := ExecCommand(m.Command)
	if err != nil {
		return err
	}
	mv.Value = value

	hostname, err := g.Hostname()
	if err != nil {
		log.Println("Getting hostname with error:", err)
	}
	mv.Endpoint = hostname
	log.Println(mv)
	mvs = append(mvs, mv)
	g.SendToTransfer(mvs)
	return nil
}

func SyncUserDefinedMetrics() {
	if !g.Config().Heartbeat.Enabled {
		return
	}

	if g.Config().Heartbeat.Addr == "" {
		return
	}

	go syncAddedMetrics()
	go syncRemovedMetrics()
}

func syncAddedMetrics() {

	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second

	for {
		time.Sleep(duration)

		hostname, err := g.Hostname()
		if err != nil {
			continue
		}

		req := model.AgentHeartbeatRequest{
			Hostname: hostname,
		}

		var resp model.AddedMetricsResponse
		err = g.HbsClient.Call("Agent.AddedMetrics", req, &resp)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		if g.Config().Debug {
			log.Println(&resp)
		}

		for _, metric := range resp.Metrics {
			ticktock.Schedule(metric.Name, &UserDefinedMetric{metric.Name, metric.Command,
			metric.Step, metric.MetricType}, &t.When{Every: t.Every(1).Minutes()})
		}
	}
}


func syncRemovedMetrics() {

	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second

	for {
		time.Sleep(duration)

		hostname, err := g.Hostname()
		if err != nil {
			continue
		}

		req := model.AgentHeartbeatRequest{
			Hostname: hostname,
		}

		var resp model.RemovedMetricsResponse
		err = g.HbsClient.Call("Agent.RemovedMetrics", req, &resp)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		if g.Config().Debug {
			log.Println(&resp)
		}

		for _, metric := range resp.Metrics {
			ticktock.Cancel(metric.Name)
		}
	}
}

