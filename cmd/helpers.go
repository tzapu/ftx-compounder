package cmd

import (
	"github.com/gin-gonic/gin"
	formatter "github.com/lacasian/logrus-module-formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initLogging() {
	logging := viper.GetString("logging")

	if verbose {
		logging = "*=debug"
	}

	if vverbose {
		logging = "*=trace"
	}

	if logging == "" {
		logging = "*=info"
	}
	viper.Set("logging", logging)

	gin.SetMode(gin.DebugMode)

	modules := formatter.NewModulesMap(logging)
	if level, exists := modules["gin"]; exists {
		if level < logrus.DebugLevel {
			gin.SetMode(gin.ReleaseMode)
		}
	} else {
		level := modules["*"]
		if level < logrus.DebugLevel {
			gin.SetMode(gin.ReleaseMode)
		}
	}

	f, err := formatter.New(modules)
	if err != nil {
		panic(err)
	}

	logrus.SetFormatter(f)

	log.Debug("Debug mode")
}
