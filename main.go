package main

import (
	"car-go/schema"
	"car-go/util"

	"github.com/gin-gonic/gin"
)

func main() {
	util.ReadConfig()
	engine := gin.Default()
	schema.LinkDb()
	engine.Run()
}
