package user

import (
	"github.com/gin-gonic/gin"
	"github.com/qiwen698/miniblog/internal/pkg/log"
	"github.com/qiwen698/miniblog/pkg/core"
)

func (ctrl *UserController) Get(c *gin.Context) {
	log.C(c).Infow("Get user function called")
	user, err := ctrl.b.Users().Get(c, c.Param("name"))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, user)
}
