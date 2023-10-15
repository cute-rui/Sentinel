package auth

import (
	"Sentinel/utils/config"
	"github.com/gin-gonic/gin"
	"time"
)

var Requested = map[string]time.Time{}

func SendEmailVerify(c *gin.Context) {
	ip := c.ClientIP()
	if ip == `` {
		c.AbortWithStatus(403)
		return
	}

	if t, ok := Requested[ip]; ok && t.Unix() > time.Now().Unix()-config.Conf.GetInt64(`Verify.Interval`) {
		c.AbortWithStatusJSON(403, gin.H{
			`Msg`: `too many requests`,
		})
		return
	}

	/*err := sendEmail()
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	c.JSON(200, gin.H{
		`Msg`: `OK`,
	})*/
}
