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

package metric

import (
	"github.com/open-falcon/falcon-plus/modules/api/config"
	"net/http"
	"github.com/gin-gonic/gin"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	db = config.Con()
	metricr := r.Group("/api/v1")
	metricr.GET("/metric", GetMetrics)
	metricr.POST("/metric", CreateMetric)
	metricr.DELETE("/metric/:metric_id", DeleteMetric)
	metricr.GET("/metric/alias", GetMetricNameAndAlia)
	metricr.POST("/namespace", CreateNameSpace)
	metricr.GET("/namespace", GetNameSpaces)
	metricr.DELETE("/namespace/:namespace_id", DeleteNameSpace)
	metricr.POST("bind/metric-host", BindMetricToHosts)
	metricr.POST("unbind/metric-host", UnBindMetricToHosts)
	metricr.POST("relation/metric-host", GetMetricBindHosts)
}
