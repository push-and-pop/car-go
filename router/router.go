package router

import (
	"car-go/controller/login"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(e *gin.Engine) {
	group := e.Group("/api/v1")
	{
		initLogin(group)
	}
}

func initLogin(r *gin.RouterGroup) {
	r.POST("/login", login.Login)
}
