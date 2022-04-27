package pack

import (
	"car-go/controller/order"
	. "car-go/schema"
	"car-go/schema/model"
	"car-go/util/json"

	"github.com/gin-gonic/gin"
)

type EnterParkReq struct {
	Location string `json:"location"`
	Number   int32  `json:"number"`
}

//直接入库，形成持续订单，并在出库后形成完整订单，并需要支付订单
func EnterPark(c *gin.Context) {
	req := &EnterParkReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	phone := c.GetString("phone")
	user := model.User{}
	err = Db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}

	park := model.CarPark{}
	err = Db.Where("location = ? and number = ?", req.Location, req.Number).First(&park).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	var interval model.TimeInterval
	err = json.Unmarshal([]byte(park.TimeInterval), &interval)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	if park.ParkState == Useing {
		c.JSON(400, gin.H{
			"err": "park is using",
		})
		return
	}
	park.ParkState = Useing
	err = Db.Save(&park).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	order := model.Order{
		Type:   order.Continue,
		State:  order.Going,
		PackId: park.ID,
	}
	err = Db.Create(&order).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "入库成功",
	})
}

type ReserveParkReq struct {
}

//预定车位，形成完整订单，并需要支付订单
func ReservePark(c *gin.Context) {

}
