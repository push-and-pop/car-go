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

type LoginRsp struct {
	Token string `json:"token"`
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
	err = Db.Where("user_name = ?", req.Phone).First(&user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	if user.Password != req.VerifyCode {
		c.JSON(400, gin.H{
			"err": "密码错误",
		})
		return
	}
	token, err := token.GenToken(user.UserName)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"token":     token,
		"user_info": user,
		"code":      200,
	})
}

type RegisterReq struct {
	Phone      string `json:"phone"`
	VerifyCode string `json:"verify_code"`
}

func Register(c *gin.Context) {
	req := &RegisterReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	user := model.User{}
	err = Db.Where("user_name = ?", req.Phone).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = model.User{
			Role:       util.User,
			Account:    0,
			CarState:   util.OutPark,
			PackId:     0,
			UserName:   req.Phone,
			Password:   req.VerifyCode,
			IsComplete: false,
			CarNumber:  "",
		}
		err = Db.Create(&user).Error
		if err != nil {
			c.JSON(400, gin.H{
				"err": err,
			})
			return
		}
		c.JSON(200, gin.H{
			"msg":  "注册成功",
			"code": 200,
		})
		return
	} else if err != nil {
		if err != nil {
			c.JSON(400, gin.H{
				"err": err.Error(),
			})
			return
		}
	}
	c.JSON(400, gin.H{
		"err": "账号已存在",
	})

}

func GetFrontPage(c *gin.Context) {
	announce := model.Announce{}
	err := Db.Last(&announce).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	park := []model.CarPark{}
	err = Db.Find(&park).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	lenpark := len(park)
	announceMsg := announce.Msg
	c.JSON(200, gin.H{
		"data": gin.H{
			"total_park": lenpark,
			"announce":   announceMsg,
		},
		"code": 200,
	})
}
