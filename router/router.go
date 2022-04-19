package router

import (
	"car-go/controller/announcement"
	"car-go/controller/login"
	"car-go/controller/message"
	"car-go/controller/pack"
	"car-go/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(e *gin.Engine) {
	e.Use(middleware.Cors())
	group := e.Group("/api")
	{
		//group.Use(middleware.Cors())
		initLogin(group)
		_ = middleware.BeforeRoute()
		//group.Use(middleware.BeforeRoute())
		initMessage(group)
		initAnnounce(group)
		initPark(group)
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
	r.GET("/annocement/get", announcement.GetAllAnnouncement)
}

//车位
func initPark(r *gin.RouterGroup) {
	r.POST("/park/creat", pack.CreatCarPark)
}
