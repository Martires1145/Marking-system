package util

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	sqlxadapter "github.com/memwey/casbin-sqlx-adapter"
	"marking/common"
	"sync"
)

type CasbinService struct {
}

var CasbinServiceApp = new(CasbinService)

// 持久化到数据库
var (
	syncedEnforcer *casbin.SyncedEnforcer
	once           sync.Once
)

func (c *CasbinService) Casbin() *casbin.SyncedEnforcer {
	once.Do(func() {
		opt := &sqlxadapter.AdapterOptions{
			DriverName:     "mysql",
			DataSourceName: common.Dsn,
			TableName:      "casbin_rule",
		}
		a := sqlxadapter.NewAdapterFromOptions(opt)
		enforcer, err := casbin.NewSyncedEnforcer("./config/model.conf", a)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		syncedEnforcer = enforcer
	})
	_ = syncedEnforcer.LoadPolicy()
	return syncedEnforcer
}

func AddGroupM(sub, g string) error {
	CasbinServiceApp.Casbin()
	_, err := syncedEnforcer.AddNamedGroupingPolicy("g", sub, g)
	return err
}

func SetRoot(uid string) error {
	return AddGroupM(uid, "root")
}
