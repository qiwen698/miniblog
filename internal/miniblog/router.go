package miniblog

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/qiwen698/miniblog/internal/miniblog/controller/v1/user"
	"github.com/qiwen698/miniblog/internal/miniblog/store"
	"github.com/qiwen698/miniblog/internal/pkg/log"
	mw "github.com/qiwen698/miniblog/internal/pkg/middleware"
	"github.com/qiwen698/miniblog/pkg/auth"
	"github.com/qiwen698/miniblog/pkg/core"
	"github.com/qiwen698/miniblog/pkg/errno"
)

func installRouters(g *gin.Engine) error {
	//注册 404 Handler
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})
	//注册 /healthz handler.
	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")
		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})
	// 注册 pprof 路由
	pprof.Register(g)

	authz, err := auth.NewAuthz(store.S.DB())
	if err != nil {
		return err
	}
	uc := user.New(store.S, authz)
	g.POST("/login", uc.Login)
	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 users 路由分组
		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)
			userv1.PUT(":name/change-password", uc.ChangePassword)
			userv1.Use(mw.Authn(), mw.Authz(authz))
			userv1.GET(":name", uc.Get)    //获取用户详情
			userv1.PUT(":name", uc.Update) //更新用户
			userv1.GET("", uc.List)        // 列出所有用户列表，只有root用户才能访问

		}
	}
	return nil
}
