package announcement

import (
	. "car-go/schema"
	"car-go/schema/model"

	"github.com/gin-gonic/gin"
)

type AnnouncementReq struct {
	Msg string `json:"msg"`
}

func PubAnnouncement(c *gin.Context) {
	req := &AnnouncementReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	err = Db.Create(&model.Announce{Msg: req.Msg}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
	})
}

type GetAllAnnouncementReq struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func GetAllAnnouncement(c *gin.Context) {
	req := &GetAllAnnouncementReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	announce := &model.Announce{}
	err = Db.Limit(int(req.Limit)).Offset(int(req.Offset)).Find(announce).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": announce,
	})
}
