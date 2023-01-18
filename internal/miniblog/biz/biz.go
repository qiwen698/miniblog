package biz

import (
	"github.com/qiwen698/miniblog/internal/miniblog/biz/user"
	"github.com/qiwen698/miniblog/internal/miniblog/store"
)

// IBiz 定义了Biz 层需要实现的方法
type IBiz interface {
	Users() user.UserBiz
}

// 确保 biz 实现了IBiz 接口.
var _ IBiz = (*biz)(nil)

type biz struct {
	ds store.IStore
}

func NewBiz(ds store.IStore) *biz {
	return &biz{
		ds: ds,
	}
}

func (b *biz) Users() user.UserBiz {
	return user.New(b.ds)
}
