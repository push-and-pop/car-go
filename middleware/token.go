package middleware

import (
	"car-go/util/token"
	"fmt"

	"github.com/gin-gonic/gin"
)

//token鉴权
func BeforeRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		claim, err := token.ParseToken(tokenStr)
		if err != nil {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  err,
			})
			c.Abort()
			return
		}
		c.Set("phone", claim.Issuer)
		fmt.Println(claim)
		c.Next()
	}
}
