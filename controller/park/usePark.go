package park

import (
	"car-go/controller/order"
	. "car-go/schema"
	"car-go/schema/model"
	"car-go/util"
	"car-go/util/json"
	"errors"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	userName := c.GetString("userName")
	user := model.User{}
	err = Db.Where("user_name = ?", userName).First(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	user.CarState = util.InPark

	park := model.CarPark{}
	err = Db.Where("location = ? and number = ?", req.Location, req.Number).First(&park).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	user.PackId = park.ID
	var isReserve bool
	if user.ReserveParkId == park.ID {
		var interval model.TimeInterval
		now := time.Now().Unix()
		err = json.Unmarshal([]byte(park.TimeInterval), &interval)
		if err != nil {
			c.JSON(400, gin.H{
				"err": err.Error(),
			})
			return
		}
		for _, value := range interval {
			if now >= value.StartTime && now <= value.EndTime {
				isReserve = true
			}
		}
	}
	if park.ParkState == Useing {
		c.JSON(400, gin.H{
			"err": "park is using",
		})
		return
	}
	user.EnterAt = time.Now().Unix()
	err = Db.Save(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	if !isReserve { //如果不是预约
		order := model.Order{
			Type:    order.Continue,
			State:   order.Going,
			UserId:  user.ID,
			PackId:  park.ID,
			StartAt: time.Now().Unix(),
		}
		err = Db.Create(&order).Error
		if err != nil {
			c.JSON(400, gin.H{
				"err": err,
			})
			return
		}
		park.OrderId = order.ID
	}
	park.ParkState = Useing
	err = Db.Save(&park).Error
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
	userName := c.GetString("userName")
	user := model.User{}
	err = tx.Where("user_name = ?", userName).First(&user).Error
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
	if user.PackId != park.ID {
		c.JSON(400, gin.H{
			"err": "此车位不是您的车辆",
		})
		return
	}
	now := time.Now().Unix()
	var extraPrice int64

	if user.ReserveParkId == park.ID {
		//取出时间段
		var interval model.TimeInterval
		err = json.Unmarshal([]byte(park.TimeInterval), &interval)
		if err != nil {
			c.JSON(400, gin.H{
				"err": err.Error(),
			})
			return
		}
		for _, value := range interval {
			if user.EnterAt >= value.StartTime && now <= value.EndTime { //在时间段内出库
				user.CarState = util.OutPark
				user.PackId = 0
				err = tx.Save(&user).Error
				if err != nil {
					tx.Rollback()
					c.JSON(400, gin.H{
						"err": err,
					})
					return
				}
				park.ParkState = Open
				park.OrderId = 0
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
				return
			}
			if user.EnterAt >= value.StartTime && now >= value.EndTime { //超过时间段
				extraPrice = int64(math.Ceil(float64(now-value.EndTime)/float64(3600))) * order.Price
				if user.Account < extraPrice {
					tx.Rollback()
					c.JSON(400, gin.H{
						"err": "用户余额不足，出库失败",
					})
					return
				}
				user.Account -= extraPrice
				user.CarState = util.OutPark
				user.PackId = 0
				err = tx.Save(&user).Error
				if err != nil {
					tx.Rollback()
					c.JSON(400, gin.H{
						"err": err,
					})
					return
				}
				park.ParkState = Open
				park.OrderId = 0
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
				break
			}
		}
	}

	Order := model.Order{}
	err = tx.Where("id = ?", park.OrderId).First(&Order).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}

	Order.EndAt = now
	Order.Price = int64(math.Ceil(float64(now-Order.StartAt)/float64(3600))) * order.Price
	Order.State = order.Pay
	if user.Account < Order.Price {
		tx.Rollback()
		c.JSON(400, gin.H{
			"err": "用户余额不足，出库失败",
		})
		return
	}
	err = tx.Save(&Order).Error
	if err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	user.Account -= Order.Price
	user.CarState = util.OutPark
	user.PackId = 0
	err = tx.Save(&user).Error
	if err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	park.ParkState = Open
	park.OrderId = 0
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
	userName := c.GetString("userName")
	user := model.User{}
	tx := Db.Begin()
	err = tx.Where("user_name = ?", userName).First(&user).Error
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
		UserId:  user.ID,
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
	user.ReserveParkId = park.ID
	err = tx.Create(&user).Error
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

func GetMyPark(c *gin.Context) {
	userName := c.GetString("userName")
	user := model.User{}
	err := Db.Where("user_name = ?", userName).First(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	park := model.CarPark{}
	err = Db.Where("id = ?", user.PackId).First(&park).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	order := model.Order{}
	err = Db.Where("id = ?", park.OrderId).First(&order).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	now := time.Now().Unix()
	c.JSON(200, gin.H{
		"code":  200,
		"order": order,
		"park":  park,
		"price": int64(math.Ceil(float64(now-order.StartAt)/float64(3600))) * order.Price,
	})
}
