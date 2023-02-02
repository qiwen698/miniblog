package user

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qiwen698/miniblog/internal/pkg/log"
	v1 "github.com/qiwen698/miniblog/pkg/api/miniblog/v1"
	"github.com/qiwen698/miniblog/pkg/core"
	"github.com/qiwen698/miniblog/pkg/errno"
	pb "github.com/qiwen698/miniblog/pkg/proto/miniblog/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// List 返回用户列表，只有 root 用户才能获取用户列表.
func (ctrl *UserController) List(c *gin.Context) {
	log.C(c).Infow("List user function called")
	var r v1.ListUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	resp, err := ctrl.b.Users().List(c, r.Offset, r.Limit)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, resp)

}

// ListUser 返回用户列表，只有 root 用户才能获取用户列表

func (ctrl *UserController) ListUser(ctx context.Context, r *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	log.C(ctx).Infow("ListUser function called")
	resp, err := ctrl.b.Users().List(ctx, int(r.Offset), int(r.Limit))
	if err != nil {
		return nil, err
	}
	users := make([]*pb.UserInfo, 0, len(resp.Users))
	for _, u := range resp.Users {
		createAt, _ := time.Parse("2006-01-02 15:04:05", u.CreatedAt)
		updateAt, _ := time.Parse("2006-01-02 15:04:05", u.UpdatedAt)
		users = append(users, &pb.UserInfo{
			Username:  u.Username,
			Nickname:  u.Nickname,
			Email:     u.Email,
			Phone:     u.Phone,
			PostCount: u.PostCount,
			CreateAt:  timestamppb.New(createAt),
			UpdateAt:  timestamppb.New(updateAt),
		})
	}
	ret := &pb.ListUserResponse{
		TotalCount: resp.TotalCount,
		Users:      users,
	}
	return ret, nil
}
