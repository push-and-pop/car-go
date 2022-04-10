package router

import (
	"car-go/controller/announcement"
	"car-go/controller/login"
	"car-go/controller/message"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(e *gin.Engine) {
	group := e.Group("/api/v1")
	{
		initLogin(group)
		initMessage(group)
		initAnnounce(group)
	}
}

//登录
func initLogin(r *gin.RouterGroup) {
	r.POST("/login", login.Login)
}

//信息
func initMessage(r *gin.RouterGroup) {
	r.POST("/checkmessage/creat", message.RegisterMessage)
}

//公告
func initAnnounce(r *gin.RouterGroup) {
	r.POST("/annocement/creat", announcement.PubAnnouncement)
	r.GET("annocement/get", announcement.GetAllAnnouncement)
}
