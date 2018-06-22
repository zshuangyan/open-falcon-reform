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

package cache

// 每个agent心跳上来的时候立马更新一下数据库是没必要的
// 缓存起来，每隔一个小时写一次DB
// 提供http接口查询机器信息，排查重名机器的时候比较有用

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"sync"
)

type SafeMetrics struct {
	sync.RWMutex
	M map[string]*model.Metric
}

var Metrics = NewSafeMetrics()

func NewSafeMetrics() *SafeMetrics {
	return &SafeMetrics{M: make(map[string]*model.Metric)}
}

func (this *SafeMetrics) Put(req *model.MetricReportRequest) {
	val := &model.Metric{
		Name:        req.Name,
		Command:     req.Command,
		Step:        req.Step,
		MetricType:  req.MetricType,
	}

	this.Lock()
	this.M[req.HostName] = val
	this.Unlock()

}

func (this *SafeMetrics) Get(hostname string) (*model.Metric, bool) {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[hostname]
	return val, exists
}


