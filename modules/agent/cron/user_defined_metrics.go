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
	"github.com/open-falcon/falcon-plus/modules/hbs/cache"
)

type UserDefinedMetricResponse struct {
	Metrics   []*cache.UserDefinedMetric
}

func SyncUserDefinedMetrics() {
	if !g.Config().Heartbeat.Enabled {
		return
	}

	if g.Config().Heartbeat.Addr == "" {
		return
	}

	go syncUserDefinedMetrics()
}

func syncUserDefinedMetrics() {

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

		var resp UserDefinedMetricResponse
		err = g.HbsClient.Call("Agent.UserDefinedMetrics", req, &resp)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		if g.Config().Debug {
			log.Println(&resp)
		}

		for _, metric := range resp.Metrics {
			log.Println(metric)
		}

	}
}
