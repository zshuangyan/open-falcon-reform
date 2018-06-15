package falcon_portal
import (
	con "github.com/open-falcon/falcon-plus/modules/api/config"
)
/*
+--------------+------------------+------+-----+---------+-------+
| Field        | Type             | Null | Key | Default | Extra |
+--------------+------------------+------+-----+---------+-------+
| namespace_id | int(10) unsigned | NO   | MUL | NULL    |       |
| metric_id    | int(10) unsigned | NO   | MUL | NULL    |       |
+--------------+------------------+------+-----+---------+-------+
 */
type NameSpaceMetric struct{
	NameSpaceID         int64    `json:"namespace_id" gorm:"column:namespace_id"`
	MetricID            int64    `json:"metric_id" gorm:"column:metric_id"`
}

func (this NameSpaceMetric) TableName() string {
	return "namespace_metric"
}

func (this NameSpaceMetric) Existing() bool {
	var tNameSpaceMetric NameSpaceMetric
	db := con.Con()
	db.Falcon.Table(this.TableName()).Where("namespace_id = ? AND metric_id = ?", this.NameSpaceID, this.MetricID).Scan(&tNameSpaceMetric)
	if tNameSpaceMetric.NameSpaceID != 0 {
		return true
	} else {
		return false
	}
}