package falcon_portal
import (
	con "github.com/open-falcon/falcon-plus/modules/api/config"
)
/*
+-------+------------------+------+-----+---------+----------------+
| Field | Type             | Null | Key | Default | Extra          |
+-------+------------------+------+-----+---------+----------------+
| id    | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
| name  | varchar(255)     | NO   | UNI | NULL    |                |
+-------+------------------+------+-----+---------+----------------+
 */
type NameSpace struct{
	ID         int64    `json:"id" gorm:"column:id"`
	Name       string   `json:"name" gorm:"column:name"`
}

func (this NameSpace) TableName() string {
	return "namespace"
}

func (this NameSpace) Existing() (int64, bool) {
	db := con.Con()
	db.Falcon.Table(this.TableName()).Where("name = ?", this.Name).Scan(&this)
	if this.ID != 0 {
		return this.ID, true
	} else {
		return 0, false
	}
}