package auth

import (
	"Sentinel/dao"
	"Sentinel/dao/models"
	"Sentinel/utils/auth"
	"Sentinel/utils/config"
	"github.com/gin-gonic/gin"
)

type registerForm struct {
	User     string `form:"username" binding:"required"`
	Realname string `form:"realname" binding:"required"`
	Password string `form:"password" binding:"required"`
	Email    string `form:"email" binding:"required"`
	QQ       int64  `form:"qq" binding:"required"`
	Verify   string `form:"verify" binding:"required"`
}

func Register(c *gin.Context) {
	var form registerForm
	if c.ShouldBind(&form) != nil {
		c.AbortWithStatusJSON(400, gin.H{
			`Msg`: `Invalid Form`,
		})
		return
	}

	/*if !dao.MatchVerify(form.Email, form.Verify) {
		c.AbortWithStatusJSON(400, gin.H{
			`Msg`: `Verify Code Mismatched`,
		})

		return
	}*/

	hash, err := auth.CreateHash(form.Password, auth.DefaultArgon2IDParams)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	user := &models.User{
		Username:   form.User,
		Realname:   form.Realname,
		Hash:       hash,
		Email:      form.Email,
		QQ:         form.QQ,
		Permission: "NOT APPROVED",
		IsAdmin:    false,
	}

	//todo: check duplicated error
	err = dao.CreateUser(user)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	token, err := auth.NewJWTClaim(user).Sign(config.Conf.GetString(`JWT.AccessSecret`))

	c.JSON(200, gin.H{
		`Msg`: `OK`,
		`Data`: gin.H{
			`Token`: token,
		},
	})
}
