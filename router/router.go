package router

import (
	"car-go/controller/login"

	"github.com/gin-gonic/gin"
)

type router struct {
	*gin.Engine
}

func RegisterRouter(e *gin.Engine) {
	r := router{
		e,
	}
	group := r.Group("/api/v1")
	{

	}
	r.initLogin()
}

func (r *router) initLogin() {
	r.POST("/login", login.Login)
}
