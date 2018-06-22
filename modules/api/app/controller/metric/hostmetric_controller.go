package metric

import (
	"fmt"

	"github.com/gin-gonic/gin"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
	"github.com/open-falcon/falcon-plus/common/model"
	"net/http"
	"strings"
)

type APIBindMetricToHosts struct {
	MetricID int64 `json:"metric_id" binding:"required"`
	Hosts []int64  `json:"host_ids" binding:"required"`
}

func BindMetricToHosts(c *gin.Context) {
	var inputs APIBindMetricToHosts
	ecode := -1
	if err := c.Bind(&inputs); err != nil {
		h.JSONResponse(c, badstatus, ecode, err)
		return
	}
	metric := f.Metric{ID: inputs.MetricID}
	if dt := db.Falcon.Find(&metric); dt.Error != nil {
		h.JSONResponse(c, expecstatus, ecode, dt.Error)
		return
	}
	var errors []string
	for _, host_id := range inputs.Hosts {
		ahost := &f.Host{ID: host_id}
		if dt:= db.Falcon.Find(&ahost); dt.Error != nil{
			errors = append(errors, fmt.Sprintf("host: %v not exist\n", host_id))
		} else {
			if dt := db.Falcon.Create(&f.HostMetric{MetricID: metric.ID, HostID: ahost.ID}); dt.Error != nil {
				errors = append(errors, fmt.Sprintf("bound to host: %v failed for reason: %s\n", host_id, dt.Error.Error()))
			}
			if dt := db.Falcon.Create(&model.UserDefinedMetricHost{metric.Name, metric.Command,
			metric.Step, metric.MetricType, metric.ValueType, ahost.ID}); dt.Error != nil {
				errors = append(errors, fmt.Sprintf("bound to host: %v failed for reason: %s\n", host_id, dt.Error.Error()))
			}
		}
	}
	var error_msg string
	if len(errors) != 0 {
		error_msg = strings.Join(errors, "")
		ecode = -1

	} else {
		error_msg = "bound succeed"
		ecode = 0
	}
	h.JSONResponse(c, http.StatusOK, 0, error_msg)
	return
}

func UnBindMetricToHosts(c *gin.Context) {
	var inputs APIBindMetricToHosts
	ecode := -1
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	metric := f.Metric{ID: inputs.MetricID}
	if dt := db.Falcon.Find(&metric); dt.Error != nil {
		h.JSONResponse(c, expecstatus, ecode, dt.Error)
		return
	}
	var errors []string
	for _, host_id := range inputs.Hosts {
		ahost := &f.Host{ID: host_id}
		if dt:= db.Falcon.Find(&ahost); dt.Error != nil{
			errors = append(errors, fmt.Sprintf("host: %v not exist\n", host_id))
		} else {
			if dt := db.Falcon.Where("metric_id = ? AND host_id = ?", inputs.MetricID, host_id).Delete(&f.HostMetric{}); dt.Error != nil {
				errors = append(errors, fmt.Sprintf("host: %v unbound failed for reason: %s\n", host_id, dt.Error.Error()))
			}
		}

	}
	var error_msg string
	if len(errors) != 0 {
		error_msg = strings.Join(errors, "")
		ecode = -1

	} else {
		error_msg = "unbound succeed"
		ecode = 0
	}
	h.JSONResponse(c, http.StatusOK, 0, error_msg)
	return
}

type APIMetricBindHosts struct {
	MetricID int64  `json:"metric_id" form:"metric_id" binding:"required"`
	Bind     bool   `json:"bind" form:"bind"`
	Status   int    `json:"status" form:"status"`
	Q        string `json:"q" form:"q"`
	Limit    int    `json:"limit" form:"limit"`
	Page     int    `json:"page" form:"page"`
}

type HostIDAndName struct {
	ID   int      `json:"id" gorm:"column:id"`
	Name string   `json:"name" gorm:"column:hostname"`
}

func GetMetricBindHosts(c *gin.Context){
	hbInterval := 5
	inputs := APIMetricBindHosts{
		Status: 2,
		Limit: 50,
		Page:  1,
	}
	ecode := -1
	if err := c.Bind(&inputs); err != nil {
		h.JSONResponse(c, badstatus, ecode, err)
		return
	}
	metric := f.Metric{ID: inputs.MetricID}
	if dt := db.Falcon.Find(&metric); dt.Error != nil {
		h.JSONResponse(c, expecstatus, ecode, dt.Error)
		return
	}
	sql := "SELECT host.id, host.hostname, host.ip, date_format(hb_at, '%Y-%m-%d %T') as updated, " +
		fmt.Sprintf("CASE WHEN TIMESTAMPDIFF(minute, host.hb_at, NOW()) < %v THEN 1 ELSE 0 END AS status", hbInterval)
	if inputs.Bind {
		sql += fmt.Sprintf(" FROM host WHERE host.id %s IN (SELECT host_id FROM host_metric WHERE host_metric.metric_id = %v)", "", metric.ID)
	} else {
		sql += fmt.Sprintf(" FROM host WHERE host.id %s IN (SELECT host_id FROM host_metric WHERE host_metric.metric_id = %v)", "NOT", metric.ID)
	}
	if inputs.Q != ""{
		sql += " AND host.hostname REGEXP " + "'.*" + inputs.Q + ".*'"
	}
	if inputs.Status == 0 {
		sql += " AND " + fmt.Sprintf("TIMESTAMPDIFF(minute, host.hb_at, NOW()) > %v", hbInterval)
	} else if inputs.Status == 1 {
		sql += " AND " + fmt.Sprintf("TIMESTAMPDIFF(minute, host.hb_at, NOW()) <= %v", hbInterval)
	}
	var offset = 0
	if inputs.Page > 1 {
		offset = (inputs.Page - 1) * inputs.Limit
	}
	var hosts []f.Host
	var count int
	dt := db.Falcon.Raw(sql)
	dt.Count(&count)
	dt.Limit(inputs.Limit).Offset(offset).Scan(&hosts)
	if dt.Error != nil {
		h.JSONResponse(c, expecstatus, ecode, dt.Error)
		return
	}
	h.JSONResponse(c, http.StatusOK, 0, fmt.Sprintf("succeed get hosts bound to metric:%v", metric.ID), &CountResult{count, hosts})
	return
}