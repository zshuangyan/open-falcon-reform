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

import (
	"github.com/open-falcon/falcon-plus/modules/hbs/db"
	"github.com/open-falcon/falcon-plus/common/model"
)


func GetAddedMetrics(hostname string) ([]*model.AddedMetric, error) {
	ret := []*model.AddedMetric{}
	hosts, err := db.QueryHosts()
	if err != nil {
		return ret, nil
	}
	hid, exists := hosts[hostname]
	if !exists {
		return ret, nil
	}
	return db.QueryAddedMetrics(hid)
}

func GetRemovedMetrics(hostname string) ([]*model.RemovedMetric, error) {
	ret := []*model.RemovedMetric{}
	hosts, err := db.QueryHosts()
	if err != nil {
		return ret, nil
	}
	hid, exists := hosts[hostname]
	if !exists {
		return ret, nil
	}
	return db.QueryRemovedMetrics(hid)
}