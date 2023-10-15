package auth

import (
	"Sentinel/dao"
	"Sentinel/dao/models"
	"Sentinel/utils/auth"
	"Sentinel/utils/config"
	"github.com/gin-gonic/gin"
)

type loginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
	Type     string `form:"type" binding:"required"`
	//MFA
	//CAPTCHA
	Remember bool `form:"remember"`
}

func Login(c *gin.Context) {
	var form loginForm
	if c.ShouldBind(&form) != nil {
		c.AbortWithStatusJSON(400, gin.H{
			`Msg`: `Invalid Form`,
		})

		return
	}

	var (
		has  bool
		user *models.User
		err  error
	)

	switch form.Type {
	case `EMAIL`:
		has, user, err = dao.FindUser(form.User, ``)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		if !has {
			c.AbortWithStatusJSON(400, gin.H{
				`Msg`: `User Not Found`,
			})
			return
		}

		if user.Permission == `NOT APPROVED` {
			c.AbortWithStatusJSON(403, gin.H{
				`Msg`: `User Not Activated Yet`,
			})
			return
		}

	case `USERNAME`:
		fallthrough
	default:
		has, user, err = dao.FindUser(``, form.User)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		if !has {
			c.AbortWithStatusJSON(400, gin.H{
				`Msg`: `User Not Found`,
			})
			return
		}

		if user.Permission == `NOT APPROVED` {
			c.AbortWithStatusJSON(403, gin.H{
				`Msg`: `User Not Activated Yet`,
			})
			return
		}
	}

	if ok, err := auth.ComparePasswordAndHash(form.Password, user.Hash); !ok || err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			`Msg`: `Password Mismatch`,
		})
		return
	}

	token, err := auth.NewJWTClaim(user, auth.WithRemember(form.Remember)).Sign(config.Conf.GetString(`JWT.AccessSecret`))
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	c.JSON(200, gin.H{
		`Msg`: `OK`,
		`Data`: gin.H{
			`Token`: token,
		},
	})
}
