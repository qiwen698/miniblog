package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/qiwen698/miniblog/internal/pkg/log"
	v1 "github.com/qiwen698/miniblog/pkg/api/miniblog/v1"
	"github.com/qiwen698/miniblog/pkg/core"
	"github.com/qiwen698/miniblog/pkg/errno"
)

func (ctrl *UserController) ChangePassword(c *gin.Context) {
	log.C(c).Infow("Change password user function called")
	var r v1.ChangePasswordRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)
		return
	}
	if err := ctrl.b.Users().ChangePassword(c, c.Param("name"), &r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, nil)
}
