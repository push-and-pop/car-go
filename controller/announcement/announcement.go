package announcement

import (
	. "car-go/schema"

	"github.com/gin-gonic/gin"
)

type Announcement struct {
	Msg string `json:"msg"`
}

func PubAnnouncement(c *gin.Context) {
	req := &Announcement{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
	Db.Model()
}
