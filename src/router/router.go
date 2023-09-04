package router

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tminaorg/brzaguza/src/config"
)

func Setup(config *config.Config) {
	router := gin.Default()

	SetupSearch(config, router)

	router.Run(fmt.Sprintf(":%v", strconv.Itoa(config.Server.Port)))
}
