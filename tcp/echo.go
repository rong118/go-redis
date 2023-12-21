package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/rong118/go_mini_redis/lib/logger"
	"github.com/rong118/go_mini_redis/lib/sync/atomic"
	"github.com/rong118/go_mini_redis/lib/sync/wait"
)

type EchoClient struct {
  Conn net.Conn
  Waiting wait.Wait
}

func (e *EchoClient) Close() error {
  e.Waiting.WaitWithTimeout(10 * time.Second)
  err := e.Conn.Close()
  return err
}

type EchoHandler struct {
  activeConn sync.Map
  closing atomic.Boolean
}

func MakeHandler() *EchoHandler{
  return &EchoHandler{}
}

func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn){
  if handler.closing.Get() {
    _ = conn.Close()
  }

  // Create internal client
  client := &EchoClient{
    Conn: conn,
  }

  handler.activeConn.Store(client, struct{}{})
  
  reader := bufio.NewReader(conn)

  for{
    msg, err := reader.ReadString('\n')

    // Error handling
    if err != nil {
      if err == io.EOF {
        logger.Info("Connecting close")
        handler.activeConn.Delete(client)
      }else{
        logger.Warn(err)
      }
      return
    }

    // Write msg back to client
    client.Waiting.Add(1)
    b := []byte(msg)
    _, _ = conn.Write(b)
    client.Waiting.Done()
  }
}

func (handler *EchoHandler) Close() error {
  logger.Info("handler close")
  handler.closing.Set(true)
  
  handler.activeConn.Range(func(key, value interface{}) bool {
    client := key.(*EchoClient)
    _ = client.Conn.Close()
    return true
  })

  return nil
}

