package metric

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
)

func GetMetrics(c *gin.Context) {
	var (
		limit int
		page  int
		err   error
	)
	pageTmp := c.DefaultQuery("page", "")
	limitTmp := c.DefaultQuery("limit", "")
	q := c.DefaultQuery("q", ".+")
	page, limit, err = h.PageParser(pageTmp, limitTmp)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	var metrics []f.Metric
	var dt *gorm.DB
	if limit != -1 && page != -1 {
		dt = db.Falcon.Raw(fmt.Sprintf("SELECT * from metric where name regexp '%s' limit %d,%d", q, page, limit)).Scan(&metrics)
	} else {
		dt = db.Falcon.Table("metric").Where("name regexp ?", q).Find(&metrics)
	}
	if dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, metrics)
	return
}