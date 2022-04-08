package login

import (
	. "car-go/schema"
	"car-go/schema/model"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginReq struct {
	Phone      string `json:"phone"`
	VerifyCode string `json:"verify_code"`
}

func Login(c *gin.Context) {
	req := &LoginReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(401, gin.H{
			"err": err,
		})
	}
	user := model.User{}
	err = Db.Where("phone = ?", req.Phone).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Db.Create(&model.User{
			Phone: req.Phone,
		})
	}
	c.JSON(200, gin.H{
		"code": 200,
	})
}
