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

package db

import (
	"github.com/open-falcon/falcon-plus/common/model"
	"fmt"
	"log"
)

func QueryAddedMetrics(hid int) ([]*model.AddedMetric, error) {
	var m []*model.AddedMetric

	sql := fmt.Sprintf("select name, command, step, metric_type from user_defined_metric where host_id=%v and status=0", hid)
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return m, err
	}

	defer rows.Close()
	for rows.Next() {
		var (
			name  string
			command string
			step int64
			metric_type string
		)

		err = rows.Scan(&name, &command, &step, &metric_type)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		m = append(m, &model.AddedMetric{name, command, step, metric_type})

		updateSql := fmt.Sprintf("update user_defined_metric set status=1 where name='%s' and host_id=%v", name, hid)
		_, err := DB.Query(updateSql)
		if err != nil {
			log.Println("ERROR:", err)
		}
	}

	return m, nil
}

func QueryRemovedMetrics(hid int) ([]*model.RemovedMetric, error) {
	var m []*model.RemovedMetric
	sql := fmt.Sprintf("select name from user_defined_metric where host_id=%v and status=2", hid)
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return m, err
	}

	defer rows.Close()
	for rows.Next() {
		var name string

		err = rows.Scan(&name)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		m = append(m, &model.RemovedMetric{name})

		updateSql := fmt.Sprintf("update user_defined_metric set status=3 where name='%s' and host_id=%v", name, hid)
		_, err := DB.Query(updateSql)
		if err != nil {
			log.Println("ERROR:", err)
		}
	}

	return m, nil
}

