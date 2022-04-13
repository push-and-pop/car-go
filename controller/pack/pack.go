package pack

import "github.com/gin-gonic/gin"

func GetParkList(c *gin.Context) {

}

func CreatCarPark(c *gin.Context) {

	c.JSON(200, gin.H{
		"code": 200,
	})
}
