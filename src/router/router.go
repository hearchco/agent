package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tminaorg/brzaguza/src/config"
)

var router = gin.Default()

func Setup(config *config.Config) {
	Search(config)
	router.Run(fmt.Sprintf(":%v", config.Server.Port))
}
