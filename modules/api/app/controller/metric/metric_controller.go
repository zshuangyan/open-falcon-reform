package metric

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
	"fmt"
	"net/http"
	"strconv"
)


type APIMetricRegexpQueryInputs struct {
	Q         string `json:"q" form:"q"`
	Limit     int    `json:"limit" form:"limit"`
	Page      int    `json:"page" form:"page"`
	NMSID     int64  `json:"namespace_id" form:"namespace_id"`
}

func GetMetrics(c *gin.Context){
	inputs := APIMetricRegexpQueryInputs{
		//set default is 500
		Limit: 50,
		Page:  1,
	}
	ecode := -1
	if err := c.Bind(&inputs); err != nil {
		h.JSONResponse(c, badstatus, ecode, err)
		return
	}

	offset := 0
	if inputs.Page > 1 {
		offset = (inputs.Page - 1) * inputs.Limit
	}

	var dt *gorm.DB
	dt = db.Falcon.Table("metric").Select("*")
	if inputs.NMSID != 0 {
		dt = dt.Joins("JOIN namespace_metric ON namespace_metric.metric_id = metric.id").Where("namespace_metric.namespace_id = ?", inputs.NMSID)
	}
	if inputs.Q != "" {
		q := ".*" + inputs.Q + ".*"
		dt = dt.Where("name regexp ?", q)
	}

	var metrics []f.Metric
	dt.Limit(inputs.Limit).Offset(offset).Scan(&metrics)
	if dt.Error != nil {
		h.JSONResponse(c, http.StatusBadRequest, ecode, dt.Error)
		return
	}

	h.JSONResponse(c, http.StatusOK, 0, "get metrics succeed", metrics)
	return
}



func GetMetricNameAndAlia(c *gin.Context) {
	var results []f.MetricNameAndAlias
	var dt *gorm.DB
	ecode := -1
	dt = db.Falcon.Raw(fmt.Sprintf("SELECT name, alias from metric")).Scan(&results)
	if dt.Error != nil {
		h.JSONR(c, expecstatus, ecode, dt.Error)
		return
	}
	h.JSONR(c, http.StatusOK, 0, "get metric alias succedd", results)
	return
}

type APICreateMetric struct {
	Name        string `json:"name" binding:"required"`
	NMSName     string `json:"namespace_name"`
	NMSID       int64  `json:"namespace_id"`
	Alias       string `json:"alias"`
	Command     string `json:"command" binding:"required"`
	Step        int    `json:"step"  binding:"required"`
	MetricType  string `json:"metric_type"  binding:"required"`
	ValueType   string `json:"value_type"  binding:"required"`
	Unit        string `json:"unit"`
}

func CreateMetric(c *gin.Context) {
	var inputs APICreateMetric
	ecode := -1
	if err := c.Bind(&inputs); err != nil {
		h.JSONResponse(c, badstatus, ecode, err)
		return
	}
	var namespace f.NameSpace
	if inputs.NMSID != 0 && inputs.NMSName != ""{
		h.JSONResponse(c, expecstatus, ecode,"cannot set namespace_id and namespace_name at the same time!")
	}
	if inputs.NMSID != 0{
		namespace = f.NameSpace{ID: inputs.NMSID}
		if dt := db.Falcon.Find(&namespace); dt.Error != nil {
			h.JSONResponse(c, expecstatus, ecode, dt.Error)
			return
		}
	}
	tx := db.Falcon.Begin()
	if inputs.NMSName != ""{
		namespace = f.NameSpace{Name: inputs.NMSName}
		if id, ok := namespace.Existing(); ok {
			namespace.ID = id
		} else if dt := tx.Create(&namespace); dt.Error != nil {
			tx.Rollback()
			h.JSONResponse(c, expecstatus, ecode, dt.Error)
			return
		}

	}
	name := namespace.Name + "." + inputs.Name
	metric := f.Metric{Name: name, Alias:inputs.Alias, Command:inputs.Command, Step:inputs.Step,
	MetricType:inputs.MetricType, ValueType:inputs.ValueType, Unit:inputs.Unit}
	if dt := tx.Create(&metric); dt.Error != nil {
		tx.Rollback()
		h.JSONResponse(c, expecstatus, ecode, dt.Error)
		return
	}
	if dt := tx.Create(&f.NameSpaceMetric{NameSpaceID: namespace.ID, MetricID: metric.ID}); dt.Error != nil {
		tx.Rollback()
		h.JSONResponse(c, expecstatus, ecode, dt.Error)
		return
	}
	tx.Commit()
	h.JSONResponse(c, http.StatusOK, 0, "success")
	return
}

func DeleteMetric(c *gin.Context) {
	MetricIDTmp := c.Params.ByName("metric_id")
	ecode := -1
	if MetricIDTmp == "" {
		h.JSONResponse(c, badstatus, ecode, "metric id is missing")
		return
	}
	MetricID, err := strconv.Atoi(MetricIDTmp)
	if err != nil {
		h.JSONResponse(c, badstatus, ecode, "metric id should be int")
		return
	}
	tx := db.Falcon.Begin()
	var host_ids []int64
	//check metric has not bound to any hosts
	db.Falcon.Table("host_metric").Select("*").Joins(
		"JOIN metric ON metric.id = host_metric.metric_id").Where("metric.id = ?", MetricID).Find(&host_ids)
	if len(host_ids) != 0 {
		h.JSONResponse(c, badstatus, ecode, fmt.Sprintf("metric: %v has bound to hosts: %v", MetricID, host_ids))
		return
	}
	//delete metric referance of namespace_metric table
	if dt := tx.Where("metric_id = ?", MetricID).Delete(&f.NameSpaceMetric{}); dt.Error != nil {
		h.JSONResponse(c, expecstatus, ecode, fmt.Sprintf("delete namespace_metric got error: %v", dt.Error))
		tx.Rollback()
		return
	}
	//finally delete metric
	if dt := tx.Delete(&f.Metric{ID: int64(MetricID)}); dt.Error != nil {
		h.JSONResponse(c, expecstatus, ecode, dt.Error)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONResponse(c, http.StatusOK, 0, fmt.Sprintf("metric: %v has been deleted", MetricID))
	return
}