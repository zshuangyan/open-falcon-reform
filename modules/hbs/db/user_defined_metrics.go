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
	"fmt"
	"log"
	"github.com/open-falcon/falcon-plus/modules/hbs/cache"
)

func QueryUserDefinedMetrics(hid int) ([]*cache.UserDefinedMetric, error) {
	var m []*cache.UserDefinedMetric

	sql := fmt.Sprint("select id, metric_name, command from user_defined_metrics where host_id='%s' and status=0", hid)
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return m, err
	}

	defer rows.Close()
	for rows.Next() {
		var (
			id int
			metric_name  string
			command string
		)

		err = rows.Scan(&id, &metric_name, &command)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		m = append(m, &cache.UserDefinedMetric{metric_name, command})

		update_sql := fmt.Sprint("update user_defined_metrics set status=1 where id='%s'", id)
		_, err := DB.Query(update_sql)
		if err != nil {
			log.Println("ERROR:", err)
		}
	}

	return m, nil
}

