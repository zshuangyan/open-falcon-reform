package metric

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"fmt"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
	"net/http"
	"strings"
	"strconv"
)

func GetNameSpaces(c *gin.Context) {
	var (
		limit int
		page  int
		err   error
	)
	pageTmp := c.DefaultQuery("page", "")
	limitTmp := c.DefaultQuery("limit", "")
	ecode := -1
	q := strings.TrimSpace(c.Query("q"))
	if q == ""{
		q = ".+"
	} else {
		q = ".*" + q + ".*"
	}
	page, limit, err = h.PageParser(pageTmp, limitTmp)
	if err != nil {
		h.JSONResponse(c, badstatus, ecode, err.Error())
		return
	}
	var namespaces []f.NameSpace
	var dt *gorm.DB
	if limit != -1 && page != -1 {
		dt = db.Falcon.Raw(fmt.Sprintf("SELECT * from namespace where name regexp '%s' limit %d,%d", q, page, limit)).Scan(&namespaces)
	} else {
		dt = db.Falcon.Table("namespace").Where("name regexp ?", q).Find(&namespaces)
	}
	if dt.Error != nil {
		h.JSONResponse(c, expecstatus, ecode, dt.Error)
		return
	}
	h.JSONResponse(c, http.StatusOK, 0, "get namespaces succeed", namespaces)
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