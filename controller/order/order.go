package order

import (
	. "car-go/schema"
	"car-go/schema/model"

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
