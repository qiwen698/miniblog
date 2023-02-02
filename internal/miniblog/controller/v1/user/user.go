package user

import (
	"github.com/qiwen698/miniblog/internal/miniblog/biz"
	"github.com/qiwen698/miniblog/internal/miniblog/store"
	"github.com/qiwen698/miniblog/pkg/auth"
	pb "github.com/qiwen698/miniblog/pkg/proto/miniblog/v1"
)

// UserController 是 user 模块在Controller 层的实现，用来处理用户模块的请求
type UserController struct {
	a *auth.Authz
	b biz.IBiz
	pb.UnimplementedMiniBlogServer
}

// New 创建一个 user controller
func New(ds store.IStore, a *auth.Authz) *UserController {
	return &UserController{
		a: a,
		b: biz.NewBiz(ds),
	}
}
