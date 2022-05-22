package user

import (
	. "car-go/schema"

	"car-go/schema/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetUserListReq struct {
	Offset int32 `query:"offset"`
	Limit  int32 `query:"limit"`
}

func GetUserList(c *gin.Context) {
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

type DeleteUserReq struct {
	Id uint `json:"id"`
}

func DeleteUser(c *gin.Context) {
	req := &DeleteUserReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	err = Db.Unscoped().Delete(&model.User{}, req.Id).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

type UpdateUserReq struct {
	Phone     string `jons:"phone"`
	Name      string `json:"trueName"`
	IdCard    string `json:"idCard"`
	CarNumber string `json:"carNumber"`
}

func UpdateUser(c *gin.Context) {
	req := &UpdateUserReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
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
		"msg":       "修改成功",
	})
}
