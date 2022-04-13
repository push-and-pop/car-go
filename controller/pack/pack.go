package pack

import (
	. "car-go/schema"
	"car-go/schema/model"

	"github.com/gin-gonic/gin"
)

func GetParkList(c *gin.Context) {

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
