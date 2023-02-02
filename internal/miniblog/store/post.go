package store

import (
	"context"

	"github.com/qiwen698/miniblog/internal/pkg/model"
	"gorm.io/gorm"
)

type PostStore interface {
	Create(ctx context.Context, post *model.PostM) error
	Get(ctx context.Context, username, postID string) (*model.PostM, error)
	Update(ctx context.Context, post *model.PostM) error
	List(ctx context.Context, username string, offset, limit int) (int64, []*model.PostM, error)
}

// PostStore 接口的实现
type posts struct {
	db *gorm.DB
}

// Get 根据 postID 查询指定用户的 post 数据库记录
func (p posts) Get(ctx context.Context, username, postID string) (*model.PostM, error) {
	var post model.PostM
	if err := p.db.Where("username = ? and postID = ?", username, postID).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil

}

// Update 更新一条 post 数据库记录

func (p posts) Update(ctx context.Context, post *model.PostM) error {
	return p.db.Save(post).Error
}

// List 根据 offset 和 limit 返回指用户的 post 列表

func (p posts) List(ctx context.Context, username string, offset, limit int) (count int64, ret []*model.PostM, err error) {
	err = p.db.Where("username = ?", username).Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	return
}

// Create 插入一条 post 记录
func (p posts) Create(ctx context.Context, post *model.PostM) error {
	return p.db.Create(&post).Error
}

// 确保 posts 实现了 PostStore 接口.
// var age int = 12 //声明变量并直接赋值.
var _ PostStore = (*posts)(nil)

func newPosts(db *gorm.DB) *posts {
	return &posts{db: db}
}
