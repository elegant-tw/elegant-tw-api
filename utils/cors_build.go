package utils

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func BuildCORSConfig(router *gin.Engine, cfg *Config) {
	corsConfig := cors.DefaultConfig()
	if cfg.CORSAllowAllOrigin {
		logrus.Info("Allow all origins: *")
		corsConfig.AllowAllOrigins = true
	} else {
		allowOrigins := strings.Split(cfg.CORSAllowOrigins, ",")
		logrus.Infof("Allow these origins: %+v", allowOrigins)
		corsConfig.AllowOrigins = allowOrigins
	}
	router.Use(cors.New(corsConfig))
}
