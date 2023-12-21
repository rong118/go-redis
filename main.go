package main

import (
	"fmt"
	"os"

	"github.com/rong118/go_mini_redis/config"
	"github.com/rong118/go_mini_redis/lib/logger"
	"github.com/rong118/go_mini_redis/tcp"
)

const configFile string = "redis.conf"

var defaultProperties = &config.ServerProperties{
  Bind: "0.0.0.0",
  Port: 6379,
}

func _fileExists(filename string) bool {
  info, err := os.Stat(filename)
  return err == nil && !info.IsDir()
}

func main() {
  logger.Setup(&logger.Settings{
    Path: "logs",
    Name: "godis",
    Ext:  "log",
    TimeFormat: "2006-01-02",
  })

  if _fileExists(configFile){
    config.SetupConfig(configFile)
  } else {
    config.Properties = defaultProperties
  }

  err := tcp.ListenAndServeWithSignal(
    &tcp.Config{
      Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
    }, 
    tcp.MakeHandler())

  if err != nil {
    logger.Error(err)
  }
}
