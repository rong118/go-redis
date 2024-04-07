package main

import (
	"flag"
	"fmt"
	"go_mini_redis/config"
	"go_mini_redis/lib/logger"
	"go_mini_redis/resp/handler"
	"go_mini_redis/tcp"
)

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

func main() {
	var echoTcp bool
	flag.BoolVar(&echoTcp, "echo", false, "Test echo tcp server")
	// Parse command-line arguments
	flag.Parse()

	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	config.Properties = defaultProperties

	var err error

    if echoTcp {
        /* This is TCP echo handler example */
        err = tcp.ListenAndServeWithSignal(
            &tcp.Config{
                Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
            },
            tcp.MakeHandler())
    }else{
        err = tcp.ListenAndServeWithSignal(
            &tcp.Config{
                Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
            },
            handler.MakeHandler())
    }

	if err != nil {
		logger.Error(err)
	}
}
