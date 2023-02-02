package auth

import (
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const (
	// casbin 访问控制模型
	aclModel = `[request_definition]
r = sub,obj,act

[policy_definition]
p = sub,obj,act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch(r.obj,p.obj) && regexMatch(r.act,p.act)

`
)

// Authz 定义一个授权器，提供授权功能

type Authz struct {
	*casbin.SyncedEnforcer
}

func NewAuthz(db *gorm.DB) (*Authz, error) {
	adapter, err := adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	m, _ := model.NewModelFromString(aclModel)
	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	enforcer.StartAutoLoadPolicy(5 * time.Second)
	a := &Authz{enforcer}
	return a, nil
}
func (a *Authz) Authorize(sub, obj, act string) (bool, error) {
	return a.Enforce(sub, obj, act)
}
