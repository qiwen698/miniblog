package user

import (
	"github.com/gin-gonic/gin"
	"github.com/qiwen698/miniblog/internal/pkg/log"
	v1 "github.com/qiwen698/miniblog/pkg/api/miniblog/v1"
	"github.com/qiwen698/miniblog/pkg/core"
	"github.com/qiwen698/miniblog/pkg/errno"
)

func (ctrl *UserController) Login(c *gin.Context) {
	log.C(c).Infow(" Login function called")
	var r v1.LoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}
	resp, err := ctrl.b.Users().Login(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, resp)
}
