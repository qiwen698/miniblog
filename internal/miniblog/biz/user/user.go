package user

import (
	"context"
	"regexp"

	v1 "github.com/qiwen698/miniblog/pkg/api/miniblog/v1"

	"github.com/jinzhu/copier"
	"github.com/qiwen698/miniblog/pkg/errno"

	"github.com/qiwen698/miniblog/internal/pkg/model"

	"github.com/qiwen698/miniblog/internal/miniblog/store"
)

// UserBiz 定义了 user 模块在 biz 层所实现的方法

type UserBiz interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error
}
type userBiz struct {
	ds store.IStore
}

func (b *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM model.UserM
	_ = copier.Copy(&userM, r)
	if err := b.ds.Users().Create(ctx, &userM); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}
		return err
	}
	return nil
}

// 确保 UserBiz 实现了UserBiz 接口
var _ UserBiz = (*userBiz)(nil)

func New(ds store.IStore) *userBiz {
	return &userBiz{
		ds: ds,
	}
}
