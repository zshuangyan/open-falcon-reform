package falcon_portal

import (
	con "github.com/open-falcon/falcon-plus/modules/api/config"
)
/*
+-------------+------------------+------+-----+---------+----------------+
| Field       | Type             | Null | Key | Default | Extra          |
+-------------+------------------+------+-----+---------+----------------+
| id          | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
| name        | varchar(255)     | NO   | UNI | NULL    |                |
| alias       | varchar(255)     | NO   |     |         |                |
| command     | varchar(500)     | NO   |     |         |                |
| step        | int(10) unsigned | NO   |     | 60      |                |
| metric_type | varchar(10)      | NO   |     | GAUGE   |                |
| value_type  | varchar(10)      | NO   |     | int     |                |
| unit        | varchar(50)      | NO   |     |         |                |
| built_in    | tinyint(1)       | NO   |     | 1       |                |
+-------------+------------------+------+-----+---------+----------------+
*/

type Metric struct {
	ID          int64  `json:"id" gorm:"column:id"`
	Name        string `json:"name" gorm:"column:name"`
	Alias       string `json:"alias" gorm:"column:alias"`
	Command     string `json:"command"  gorm:"command"`
	Step        int    `json:"step"  gorm:"step"`
	MetricType  string `json:"metric_type"  gorm:"metric_type"`
	ValueType   string `json:"value_type"  gorm:"column:value_type"`
	Unit        string `json:"unit" gorm:"column:unit"`
	BuiltIn     bool   `json:"built_in" gorm:"column:built_in"`
}

type MetricNameAndAlias struct {
	Name string   `json:"name" gorm:"column:name"`
	Alias string  `json:"alias" gorm:"column:alias"`
}

func (this Metric) TableName() string {
	return "metric"
}

func (this Metric) Existing() (int64, bool) {
	db := con.Con()
	db.Falcon.Table(this.TableName()).Where("name = ?", this.Name).Scan(&this)
	if this.ID != 0 {
		return this.ID, true
	} else {
		return 0, false
	}
}

