package announcement

import (
	. "car-go/schema"
	"car-go/schema/model"
	"strconv"

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
	Offset int32 `query:"offset"`
	Limit  int32 `query:"limit"`
}

type GetAllAnnouncementRsp struct {
	Announcement []struct {
		Date int64  `json:"date"`
		Msg  string `json:"msg"`
	} `json:"announcement"`
	Total int64 `json:"total"`
}

func GetAllAnnouncement(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	offset = limit * offset
	if offset < 0 || limit > 5 {
		c.JSON(400, gin.H{
			"err": "invalid params",
		})
		return
	}
	announce := []model.Announce{}
	err := Db.Offset(int(offset)).Limit(int(limit)).Find(&announce).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	var count int64
	err = Db.Model(&model.Announce{}).Count(&count).Error
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	rsp := GetAllAnnouncementRsp{
		Announcement: make([]struct {
			Date int64  `json:"date"`
			Msg  string `json:"msg"`
		}, len(announce)),
		Total: count,
	}
	for index, value := range announce {
		rsp.Announcement[index].Date = value.CreatedAt.Unix()
		rsp.Announcement[index].Msg = value.Msg
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": rsp,
	})
}
