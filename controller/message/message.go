package message

import (
	. "car-go/schema"
	"car-go/schema/model"

	"github.com/gin-gonic/gin"
)

type UserMessage struct {
	Name   string `json:"name"`
	IdCard string `json:"id_card"`
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
	phone := c.GetString("phone")
	user := model.User{}
	err = tx.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	user.Name = req.Name
	user.IdCard = req.IdCard
	err = tx.Save(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	tx.Commit()
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "上传成功",
	})
}
