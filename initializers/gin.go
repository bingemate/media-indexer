package initializers

import (
	"github.com/gin-gonic/gin"
)

func InitGinEngine(env Env) *gin.Engine {
	ginEngine := gin.Default()
	return ginEngine
}
