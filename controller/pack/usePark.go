package pack

import (
	"car-go/controller/order"
	. "car-go/schema"
	"car-go/schema/model"
	"car-go/util"
	"car-go/util/json"
	"time"

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
	tx := Db.Begin()
	err = tx.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	user.CarState = util.InPark
	err = tx.Save(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	park := model.CarPark{}
	err = tx.Where("location = ? and number = ?", req.Location, req.Number).First(&park).Error
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

	order := model.Order{
		Type:    order.Continue,
		State:   order.Going,
		UserId:  user.ID,
		PackId:  park.ID,
		StartAt: time.Now().Unix(),
	}
	err = tx.Create(&order).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	park.ParkState = Useing
	park.OrderId = order.ID
	err = tx.Save(&park).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	tx.Commit()
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "入库成功",
	})
}

type LeaveParkReq struct {
	Location string `json:"location"`
	Number   int32  `json:"number"`
}

//出库，订单变化，需要支付
func LeavePark(c *gin.Context) {
	req := &LeaveParkReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	tx := Db.Begin()
	phone := c.GetString("phone")
	user := model.User{}
	err = tx.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	park := model.CarPark{}
	err = tx.Where("location = ? and number = ?", req.Location, req.Number).First(&park).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	Order := model.Order{}
	err = tx.Where("id = ?", park.OrderId).First(&Order).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	now := time.Now().Unix()
	Order.EndAt = now
	Order.Price = (now - Order.StartAt) / 3600 * order.Price
	Order.State = order.UnPay
	err = tx.Save(&Order).Error
	if err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	user.CarState = util.OutPark
	err = tx.Save(&user).Error
	if err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	park.ParkState = Open
	err = tx.Save(&park).Error
	if err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	tx.Commit()
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "出库成功",
	})
}

type ReserveParkReq struct {
	Location  string `json:"location"`
	Number    int32  `json:"number"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

//预定车位，形成完整订单，只能预定第二天以后的车位，并需要支付订单，30分钟后未支付订单，订单消失
func ReservePark(c *gin.Context) {
	req := &ReserveParkReq{}
	err := c.ShouldBindJSON(req)
	if err != nil || req.StartTime >= req.EndTime {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	phone := c.GetString("phone")
	user := model.User{}
	tx := Db.Begin()
	err = tx.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	park := model.CarPark{}
	err = tx.Where("location = ? and number = ?", req.Location, req.Number).First(&park).Error
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
			"err": err.Error(),
		})
		return
	}
	if !util.HasNoIntersection(interval, req.StartTime, req.EndTime) {
		c.JSON(400, gin.H{
			"err": "时间段内已被预定",
		})
		return
	}
	interval[req.StartTime] = struct {
		StartTime int64 "json:\"start_time\""
		EndTime   int64 "json:\"end_time\""
	}{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	intervalByte, err := json.Marshal(&interval)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	park.TimeInterval = string(intervalByte)
	err = tx.Save(&park).Error
	if err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	order := model.Order{
		Type:    order.Whole,
		State:   order.UnPay,
		PackId:  park.ID,
		StartAt: req.StartTime,
		EndAt:   req.EndTime,
		Price:   (req.EndTime - req.StartTime) / 3600 * order.Price,
	}
	err = tx.Create(&order).Error
	if err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	tx.Commit()
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "预定成功",
	})
}
