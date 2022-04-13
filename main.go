package main

import (
	"car-go/router"
	"car-go/schema"
	"car-go/util"

	"github.com/gin-gonic/gin"
)

func main() {
	util.ReadConfig()
	engine := gin.Default()
	router.RegisterRouter(engine)
	schema.LinkDb()
	engine.Run(":9090")

}
