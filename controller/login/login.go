package login

import (
	. "car-go/schema"
	"car-go/schema/model"
	"car-go/util"
	"car-go/util/token"
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
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	user := model.User{}
	err = Db.Where("phone = ?", req.Phone).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = model.User{
			Phone:   req.Phone,
			Role:    util.User,
			Account: 0,
		}
		err = Db.Create(&user).Error
		if err != nil {
			c.JSON(400, gin.H{
				"err": err,
			})
			return
		}
	}
	token, err := token.GenToken(req.Phone)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"data": gin.H{
			"token":     token,
			"user_info": user,
		},
		"code": 200,
	})
}
