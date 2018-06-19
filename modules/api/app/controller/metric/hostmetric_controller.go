package metric

import (
	"fmt"

	"github.com/gin-gonic/gin"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
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
	Bind     bool   `json:"bind" form:"bind" binding:"required"`
	Q        string `json:"q" form:"q"`
	Limit    int    `json:"limit" form:"limit"`
	Page     int    `json:"page" form:"page"`
}

type HostIDAndName struct {
	ID   int      `json:"id" gorm:"column:id"`
	Name string   `json:"name" gorm:"column:hostname"`
}

func GetMetricBindHosts(c *gin.Context){
	inputs := APIMetricBindHosts{
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
	var hosts []HostIDAndName
	sql := "SELECT host.id, host.hostname FROM host WHERE host.id %s IN (SELECT host_id FROM host_metric WHERE host_metric.metric_id = %v)"
	if inputs.Bind  {
		sql = fmt.Sprintf(sql, "", metric.ID)
	} else {
		sql = fmt.Sprintf(sql, "NOT", metric.ID)
	}
	if inputs.Q != ""{
		sql += " AND host.hostname REGEXP " + "'.*" + inputs.Q + ".*'"
	}
	var offset = 0
	if inputs.Page > 1 {
		offset = (inputs.Page - 1) * inputs.Limit
	}
	sql += fmt.Sprintf(" LIMIT %v OFFSET %v", inputs.Limit, offset)
	if dt := db.Falcon.Raw(sql).Scan(&hosts); dt.Error != nil {
		h.JSONResponse(c, expecstatus, ecode, dt.Error)
		return
	}
	h.JSONResponse(c, http.StatusOK, 0, fmt.Sprintf("succeed get hosts bound to metric:%v", metric.ID), hosts)
	return
}