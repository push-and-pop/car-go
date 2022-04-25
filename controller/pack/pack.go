package pack

import (
	. "car-go/schema"
	"car-go/schema/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetParkListReq struct {
	Offset int32 `query:"offset"`
	Limit  int32 `query:"limit"`
}

func GetParkList(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	offset = limit * offset
	if offset < 0 || limit > 5 {
		c.JSON(400, gin.H{
			"err": "invalid params",
		})
		return
	}
	park := []model.CarPark{}
	err := Db.Offset(int(offset)).Limit(int(limit)).Find(&park).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"total": len(park),
			"park":  park,
		},
	})
}

type CreatCarParkReq struct {
	Location string `json:"location"`
	Number   int32  `json:"number"`
}

func CreatCarPark(c *gin.Context) {
	req := &CreatCarParkReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	err = Db.Create(&model.CarPark{Location: req.Location, Number: req.Number, ParkState: Open}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})

}

type DeleteParkReq struct {
	Id int64 `json:"id"`
}

func DeleteParkById(c *gin.Context) {
	req := &DeleteParkReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	err = Db.Delete(&model.CarPark{}, req.Id).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
