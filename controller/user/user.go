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
