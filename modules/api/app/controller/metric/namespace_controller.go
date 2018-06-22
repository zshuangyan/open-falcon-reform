package metric

import (
	"github.com/gin-gonic/gin"
	"fmt"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
	"net/http"
	"strings"
	"strconv"
)

type APINameSpaceRegexpQueryInputs struct {
	Q         string `json:"q" form:"q"`
	Limit     int    `json:"limit" form:"limit"`
	Page      int    `json:"page" form:"page"`
}

func GetNameSpaces(c *gin.Context) {
	inputs := APINameSpaceRegexpQueryInputs{
		//set default is 50
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
	dt := db.Falcon.Table("namespace").Select("*")
	q := strings.TrimSpace(c.Query("q"))
	if q != "" {
		q = ".*" + q + ".*"
		dt = dt.Where("name regexp ?", q)
	}
	var namespaces []f.NameSpace
    var count int
    dt.Count(&count).Limit(inputs.Limit).Offset(offset).Scan(&namespaces)
	if dt.Error != nil {
		h.JSONResponse(c, expecstatus, ecode, dt.Error)
		return
	}
	h.JSONResponse(c, http.StatusOK, 0, "get namespaces succeed", &CountResult{count, namespaces})
	return
}

type APICreateNamespace struct {
	Name string `json:"name" binding:"required"`
}

func CreateNameSpace(c *gin.Context) {
	var inputs APICreateNamespace
	ecode := -1
	if err := c.Bind(&inputs); err != nil {
		h.JSONResponse(c, badstatus, ecode, err)
		return
	}
	namespace := f.NameSpace{Name: inputs.Name}
	if dt := db.Falcon.Create(&namespace); dt.Error != nil {
		h.JSONResponse(c, expecstatus, ecode, dt.Error)
		return
	}
	msg := fmt.Sprint("create namespace succeed")
	h.JSONResponse(c, http.StatusOK, 0, msg, namespace)
	return
}

func DeleteNameSpace(c *gin.Context){
	NMSIDTmp := c.Params.ByName("namespace_id")
	ecode := -1
	if NMSIDTmp == "" {
		h.JSONResponse(c, badstatus, ecode, "namespace id is missing")
		return
	}
	NMSID, err := strconv.Atoi(NMSIDTmp)
	if err != nil {
		h.JSONResponse(c, badstatus, ecode, "namespace id should be int")
		return
	}
	tx := db.Falcon.Begin()
	var metric_ids []int64
	//check namespace has not bound to any metrics
	db.Falcon.Table("namespace_metric").Select("*").Joins(
		"JOIN namespace ON namespace.id = namespace_metric.namespace_id").Where(
			"namespace.id = ?", NMSID).Find(&metric_ids)
	if len(metric_ids) != 0 {
		h.JSONResponse(c, badstatus, ecode, fmt.Sprintf("namespace: %v has bound to metrics: %v", NMSID, metric_ids))
		return
	}
	//finally delete namespace
	if dt := tx.Delete(&f.NameSpace{ID: int64(NMSID)}); dt.Error != nil {
		h.JSONResponse(c, expecstatus, ecode, dt.Error)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONResponse(c, http.StatusOK, 0, fmt.Sprintf("namespace: %v has been deleted", NMSID))
	return
}