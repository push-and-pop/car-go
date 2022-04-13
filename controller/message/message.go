package message

import "github.com/gin-gonic/gin"

type UserMessage struct {
	Name   string `json:"name"`
	IdCard string `json:"id_card"`
}

//上传信息
func RegisterMessage(c *gin.Context) {
	req := &UserMessage{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err,
		})
		return
	}
}
