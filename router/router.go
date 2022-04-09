package router

import (
	"car-go/controller/login"
	"car-go/controller/message"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(e *gin.Engine) {
	group := e.Group("/api/v1")
	{
		initLogin(group)
		initMessage(group)
	}
}

func initLogin(r *gin.RouterGroup) {
	r.POST("/login", login.Login)
}

func initMessage(r *gin.RouterGroup) {
	r.POST("/checkmessage/creat", message.RegisterMessage)
}
