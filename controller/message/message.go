package message

import (
	. "car-go/schema"
	"car-go/schema/model"

	"github.com/gin-gonic/gin"
)

type UserMessage struct {
	Phone     string `jons:"phone"`
	Name      string `json:"trueName"`
	IdCard    string `json:"idCard"`
	CarNumber string `json:"carNumber"`
}

//上传信息
func RegisterMessage(c *gin.Context) {
	req := &UserMessage{}
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
	user.Phone = req.Phone
	user.CarNumber = req.CarNumber
	user.Name = req.Name
	user.IdCard = req.IdCard
	user.IsComplete = true
	err = tx.Save(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	tx.Commit()
	c.JSON(200, gin.H{
		"code":      200,
		"user_info": user,
		"msg":       "上传成功",
	})
}
