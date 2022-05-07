package router

import (
	"car-go/controller/announcement"
	"car-go/controller/login"
	"car-go/controller/message"
	"car-go/controller/order"
	"car-go/controller/park"
	"car-go/controller/user"
	"car-go/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(e *gin.Engine) {
	e.Use(middleware.Cors())
	group := e.Group("/api")
	{
		//group.Use(middleware.Cors())
		initLogin(group)
		group.Use(middleware.BeforeRoute())
		initMessage(group)
		initAnnounce(group)
		initPark(group)
		initUserPark(group)
		initUser(group)
		initOrder(group)
	}

}

//登录
func initLogin(r *gin.RouterGroup) {
	r.POST("/login", login.Login)
	r.GET("/frontPage", login.GetFrontPage)
	r.POST("/register", login.Register)
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
	r.POST("/park/creat", park.CreatCarPark)
	r.GET("/park/get", park.GetParkList)
	r.POST("/park/delete", park.DeleteParkById)
}

//用户车位
func initUserPark(r *gin.RouterGroup) {
	r.POST("/park/enter", park.EnterPark)
	r.POST("/park/leave", park.LeavePark)
	r.POST("/park/reserve", park.ReservePark)
	r.GET("/park/getself", park.GetMyPark)
}

func initUser(r *gin.RouterGroup) {
	r.GET("/user/get", user.GetUserList)
}

func initOrder(r *gin.RouterGroup) {
	r.POST("/order/pay", order.PayOrder)
	r.GET("/order/get", order.GetOrderList)
	r.POST("/order/recharge", order.RechargeAccount)
}
