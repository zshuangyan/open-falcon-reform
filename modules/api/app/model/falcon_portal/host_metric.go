package falcon_portal
import (
con "github.com/open-falcon/falcon-plus/modules/api/config"
)
/*
+-----------+------------------+------+-----+---------+-------+
| Field     | Type             | Null | Key | Default | Extra |
+-----------+------------------+------+-----+---------+-------+
| host_id   | int(10) unsigned | NO   | MUL | NULL    |       |
| metric_id | int(10) unsigned | NO   | MUL | NULL    |       |
+-----------+------------------+------+-----+---------+-------+
 */
type HostMetric struct{
	HostID         int64    `json:"host_id" gorm:"column:host_id"`
	MetricID       int64    `json:"metric_id" gorm:"column:metric_id"`
}

func (this HostMetric) TableName() string {
	return "host_metric"
}

func (this HostMetric) Existing() bool {
	var tHostMetric HostMetric
	db := con.Con()
	db.Falcon.Table(this.TableName()).Where("host_id = ? AND metric_id = ?", this.HostID, this.MetricID).Scan(&tHostMetric)
	if tHostMetric.HostID != 0 {
		return true
	} else {
		return false
	}
}
