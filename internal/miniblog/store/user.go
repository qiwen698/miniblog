package store

import (
	"context"

	"github.com/qiwen698/miniblog/internal/pkg/model"
	"gorm.io/gorm"
)

// UserStore 定义了 user 模块在 store 层所实现的方法
type UserStore interface {
	Create(ctx context.Context, user *model.UserM) error
	Get(ctx context.Context, username string) (*model.UserM, error)
	Update(ctx context.Context, user *model.UserM) error
	List(ctx context.Context, offset, limit int) (int64, []*model.UserM, error)
}

type users struct {
	db *gorm.DB
}

// List 根据 offset 和 limit 返回 user 列表
func (u *users) List(ctx context.Context, offset, limit int) (count int64, ret []*model.UserM, err error) {
	err = u.db.Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	return
}

// 确保 users 实现了 UserStore 接口
var _ UserStore = (*users)(nil)

func newUsers(db *gorm.DB) *users {
	return &users{
		db: db,
	}
}

// Create 插入一条user记录

func (u *users) Create(ctx context.Context, user *model.UserM) error {
	return u.db.Create(&user).Error
}

// 根据用户名查询指定 user 的数据库记录

func (u *users) Get(ctx context.Context, username string) (*model.UserM, error) {
	var user model.UserM
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil

}

// Update 更新一条 user 数据库记录.

func (u *users) Update(ctx context.Context, user *model.UserM) error {
	return u.db.Save(user).Error
}
