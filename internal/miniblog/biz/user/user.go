package user

import (
	"context"
	"errors"
	"regexp"
	"sync"

	"github.com/qiwen698/miniblog/internal/pkg/log"
	"golang.org/x/sync/errgroup"

	"github.com/qiwen698/miniblog/pkg/token"

	"github.com/qiwen698/miniblog/pkg/auth"
	"gorm.io/gorm"

	v1 "github.com/qiwen698/miniblog/pkg/api/miniblog/v1"

	"github.com/jinzhu/copier"
	"github.com/qiwen698/miniblog/pkg/errno"

	"github.com/qiwen698/miniblog/internal/pkg/model"

	"github.com/qiwen698/miniblog/internal/miniblog/store"
)

// UserBiz 定义了 user 模块在 biz 层所实现的方法

type UserBiz interface {
	ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error
	Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
	Create(ctx context.Context, r *v1.CreateUserRequest) error
	Get(ctx context.Context, username string) (*v1.GetUserResponse, error)
	List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error)
	Update(ctx context.Context, username string, r *v1.UpdateUserRequest) error
}
type userBiz struct {
	ds store.IStore
}

func (b *userBiz) List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error) {
	count, list, err := b.ds.Users().List(ctx, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list users from storage", "err", err)
		return nil, err
	}
	var m sync.Map
	eg, ctx := errgroup.WithContext(ctx)
	// 使用 goroutine 提高接口性能
	for _, item := range list {
		user := item
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				count, _, err := b.ds.Posts().List(ctx, user.Username, 0, 0)
				if err != nil {
					log.C(ctx).Errorw("Failed to list posts", "err", err)
				}
				m.Store(user.ID, &v1.UserInfo{
					Username:  user.Username,
					Nickname:  user.Nickname,
					Email:     user.Email,
					Phone:     user.Phone,
					PostCount: count,
					CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
					UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
				})
				return nil
			}
		})
	}
	if err := eg.Wait(); err != nil {
		log.C(ctx).Errorw("Failed to wait all function calls returned", "err", err)
		return nil, err
	}
	users := make([]*v1.UserInfo, 0, len(list))
	for _, item := range list {
		user, _ := m.Load(item.ID)
		users = append(users, user.(*v1.UserInfo))
	}
	log.C(ctx).Debugw("Get users from backend storage", "count", len(users))
	return &v1.ListUserResponse{
		TotalCount: count,
		Users:      users,
	}, nil
}

// Login 是 UserBiz 接口中 `Login` 方法的实现

func (b *userBiz) Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {
	// 获取登录用户的所有信息
	user, err := b.ds.Users().Get(ctx, r.Username)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}
	// 对比传入的明文密码和数据库中已加密过的密码是否匹配
	if err := auth.Compare(user.Password, r.Password); err != nil {
		return nil, errno.ErrPasswordIncorrect
	}
	// 如果匹配成功，说明登录成功，签发token并返回
	t, err := token.Sign(r.Username)
	if err != nil {
		return nil, errno.ErrSignToken
	}
	return &v1.LoginResponse{Token: t}, nil
}

//  Create 是 UserBiz 接口中 `Create` 方法的实现

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

func (b *userBiz) ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error {
	userM, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		return err
	}
	if err := auth.Compare(userM.Password, r.OldPassword); err != nil {
		return errno.ErrPasswordIncorrect
	}
	userM.Password, _ = auth.Encrypt(r.NewPassword)
	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return err
	}
	return nil
}

func (b *userBiz) Get(ctx context.Context, username string) (*v1.GetUserResponse, error) {
	user, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}
		return nil, err
	}
	var resp v1.GetUserResponse
	_ = copier.Copy(&resp, user)
	resp.CreatedAt = user.CreatedAt.Format("2006-01-01 15:04:05")
	resp.UpdatedAt = user.UpdatedAt.Format("2006-01-01 15:04:05")
	return &resp, nil
}
func (b *userBiz) Update(ctx context.Context, username string, user *v1.UpdateUserRequest) error {
	userM, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		return err
	}
	if user.Email != nil {
		userM.Email = *user.Email
	}
	if user.Nickname != nil {
		userM.Nickname = *user.Nickname
	}
	if user.Phone != nil {
		userM.Phone = *user.Phone
	}
	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return err
	}
	return nil

}
