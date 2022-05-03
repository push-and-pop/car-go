package order

import (
	. "car-go/schema"
	"car-go/schema/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RechargeAccountReq struct {
	Num int64 `json:"num"`
}

//充值帐户
func RechargeAccount(c *gin.Context) {
	req := RechargeAccountReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
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
	user.Account += req.Num
	err = tx.Save(&user).Error
	if err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	tx.Commit()
	c.JSON(200, gin.H{})
}

type PayOrderReq struct {
	Id uint `json:"ID"`
}

//支付订单
func PayOrder(c *gin.Context) {
	req := PayOrderReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
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
			"err": err.Error(),
		})
		return
	}
	Order := model.Order{}
	err = tx.Where("id = ?", req.Id).First(&Order).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	if user.Account < Order.Price {
		c.JSON(400, gin.H{
			"err": "用户余额不足",
		})
		return
	}
	user.Account -= Order.Price
	Order.State = Pay
	err = tx.Save(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	err = tx.Save(&Order).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	tx.Commit()
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "支付成功",
	})
}

type GetOrderListReq struct {
}

//获取订单列表
func GetOrderList(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	offset = limit * offset
	if offset < 0 || limit > 5 {
		c.JSON(400, gin.H{
			"err": "invalid params",
		})
		return
	}
	userName := c.GetString("userName")
	user := model.User{}
	err := Db.Where("user_name = ?", userName).First(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	orders := []model.Order{}
	err = Db.Offset(int(offset)).Limit(int(limit)).Where("user_id = ?", user.ID).Find(&orders).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"total": len(orders),
			"order": orders,
		},
	})
}
