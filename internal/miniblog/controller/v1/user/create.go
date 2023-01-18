package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/qiwen698/miniblog/internal/pkg/log"
	v1 "github.com/qiwen698/miniblog/pkg/api/miniblog/v1"
	"github.com/qiwen698/miniblog/pkg/core"
	"github.com/qiwen698/miniblog/pkg/errno"
)

func (ctrl *UserController) Create(c *gin.Context) {
	log.C(c).Infow("Create user function called")
	var r v1.CreateUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)
		return
	}
	if err := ctrl.b.Users().Create(c, &r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, nil)
}
