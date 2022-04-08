package login

import "github.com/gin-gonic/gin"

type LoginReq struct {
	Phone      string
	VerifyCode string
}

func Login(c *gin.Context) {

}
