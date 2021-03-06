// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package falcon_portal

import (
	"errors"
	"fmt"

	con "github.com/open-falcon/falcon-plus/modules/api/config"
)

type Host struct {
	ID            int64  `json:"id"             gorm:"column:id"`
	Hostname      string `json:"hostname"       gorm:"column:hostname"`
	IP            string `json:"ip"             gorm:"column:ip"`
	Updated       string `json:"updated"        gorm:"column:updated"`
	Status        int    `json:"status"         gorm:"column:status"`
}

func (this Host) TableName() string {
	return "host"
}

func (this Host) Existing() (int64, bool) {
	db := con.Con()
	db.Falcon.Table(this.TableName()).Where("hostname = ?", this.Hostname).Scan(&this)
	if this.ID != 0 {
		return this.ID, true
	} else {
		return 0, false
	}
}

func (this Host) RelatedGrp() (Grps []HostGroup) {
	db := con.Con()
	grpHost := []GrpHost{}
	db.Falcon.Select("grp_id").Where("host_id = ?", this.ID).Find(&grpHost)
	tids := []int64{}
	for _, t := range grpHost {
		tids = append(tids, t.GrpID)
	}
	tidStr, _ := arrInt64ToString(tids)
	Grps = []HostGroup{}
	db.Falcon.Where(fmt.Sprintf("id in (%s)", tidStr)).Find(&Grps)
	return
}

func (this Host) RelatedTpl() (tpls []Template) {
	db := con.Con()
	grps := this.RelatedGrp()
	gids := []int64{}
	for _, g := range grps {
		gids = append(gids, g.ID)
	}
	gidStr, _ := arrInt64ToString(gids)
	grpTpls := []GrpTpl{}
	db.Falcon.Select("tpl_id").Where(fmt.Sprintf("grp_id in (%s)", gidStr)).Find(&grpTpls)
	tids := []int64{}
	for _, t := range grpTpls {
		tids = append(tids, t.TplID)
	}
	tidStr, _ := arrInt64ToString(tids)
	tpls = []Template{}
	db.Falcon.Where(fmt.Sprintf("id in (%s)", tidStr)).Find(&tpls)
	return
}

func arrInt64ToString(arr []int64) (result string, err error) {
	result = ""
	for indx, a := range arr {
		if indx == 0 {
			result = fmt.Sprintf("%v", a)
		} else {
			result = fmt.Sprintf("%v,%v", result, a)
		}
	}
	if result == "" {
		err = errors.New(fmt.Sprintf("array is empty, err: %v", arr))
	}
	return
}
