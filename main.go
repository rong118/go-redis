package main

import (
	"fmt"
	"github.com/rong118/go_mini_redis/config"
	"github.com/rong118/go_mini_redis/lib/logger"
	"github.com/rong118/go_mini_redis/resp/handler"
	"github.com/rong118/go_mini_redis/tcp"
)

var defaultProperties = &config.ServerProperties{
  Bind: "0.0.0.0",
  Port: 6379,
}

func main() {
  logger.Setup(&logger.Settings{
    Path: "logs",
    Name: "godis",
    Ext:  "log",
    TimeFormat: "2006-01-02",
  })

  config.Properties = defaultProperties

  /* This is TCP echo handler example */ 
  // err := tcp.ListenAndServeWithSignal(
  //   &tcp.Config{
  //     Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
  //   }, 
  //   tcp.MakeHandler())

  err := tcp.ListenAndServeWithSignal(
    &tcp.Config{
      Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
    }, 
    handler.MakeHandler())
  
  if err != nil {
    logger.Error(err)
  }
}
