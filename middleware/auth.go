package middleware

import (
	"Sentinel/dao"
	"Sentinel/utils/auth"
	"Sentinel/utils/config"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware(c *gin.Context) {
	header := c.GetHeader(`Authorization`)
	if header == `` {
		c.AbortWithStatusJSON(401, gin.H{
			`Msg`: `Unauthorized`,
		})
		return
	}

	claim, err := auth.StringToJWTClaim(auth.TrimBearerScheme(header), config.Conf.GetString(`JWT.AccessSecret`))
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			`Msg`: `Token Invalid`,
		})
		return
	}

	c.Set(`identity`, claim)
	c.Next()
}

func LockMiddleware(c *gin.Context) {
	claimRaw, ok := c.Get(`identity`)
	if !ok {
		c.AbortWithStatus(403)
		return
	}

	claim, ok := claimRaw.(*auth.JWTClaim)
	if !ok {
		c.AbortWithStatus(403)
		return
	}

	if !dao.IsAdmin(claim.UserID) {
		c.AbortWithStatus(403)
		return
	}

	hidden := c.GetHeader(`Version`)
	if hidden != config.Conf.GetString(`Admin.Secret`) {
		c.AbortWithStatus(403)
		return
	}

	c.Next()
}
