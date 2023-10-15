package router

import (
	"Sentinel/middleware"
	adminSvr "Sentinel/service/admin"
	authSvr "Sentinel/service/auth"
	oauth2 "Sentinel/service/oauth"
	"Sentinel/utils/config"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()

	auth := r.Group(`/auth`)
	{
		auth.POST(`/login`, authSvr.Login)
		auth.POST(`/register`, authSvr.Register)
		//auth.POST(`/forgot`)
		//auth.POST(`/reset`)
		//auth.POST(`/captcha`)
		//auth.POST(`/verify`)
	}

	oauth := r.Group(`/oauth2`)
	oidc := oauth2.SetupOIDC()
	{
		oauth.POST(`/login`, gin.WrapF(oidc.LoginFunc()))
		oauth.GET(`/.well-known/openid-configuration`, gin.WrapF(oidc.DiscoveryHandler()))
		oauth.Any("/oidc/*w", gin.WrapH(oidc.GetMainHandler()))
	}

	admin := r.Group(`/admin`)
	admin.Use(middleware.JWTMiddleware, middleware.LockMiddleware)
	{
		admin.GET(`/unapproved`, adminSvr.GetUnapproved)
		admin.POST(`/approve`, adminSvr.Approve)
	}

	if err := r.Run(config.Conf.GetString(`HTTP.ListenAddr`)); err != nil {
		panic(err)
	}
}
