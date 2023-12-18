package tcp

/**
 * A tcp server
 */

import (
	"context"
	"net"
  "github.com/rong118/go_mini_redis/interface/tcp"
)

// Config stores tcp server properties
type Config struct {
	Address    string        `yaml:"address"`
}

// ClientCounter Record the number of clients in the current Godis server
var ClientCounter int

// ListenAndServeWithSignal binds port and handle requests, blocking until receive stop signal
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
  listerner, err := net.Listen("tcp", cfg.Address)
  if err != nil {
    return err
  }

  closeChan := make(chan struct{})

  ListenAndServe(listerner, handler, closeChan)


  return nil
}

// ListenAndServe binds port and handle requests, blocking until close
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
  ctx := context.Background()
  for true {
    conn, err := listener.Accept()
    if err != nil {
      break
    }

    go func (){
      handler.Handle(ctx, conn)
    }()
  }
}
