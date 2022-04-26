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
	phone := c.GetString("puhone")
	user := model.User{}
	tx := Db.Begin()
	err = tx.Where("phone = ?", phone).First(&user).Error
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
}

//支付订单
func PayOrder(c *gin.Context) {

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
	user := []model.User{}
	err := Db.Offset(int(offset)).Limit(int(limit)).Find(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"total": len(user),
			"user":  user,
		},
	})
}
