package handler

import (
	"context"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/rong118/go_mini_redis/database"
	dataBaseFace "github.com/rong118/go_mini_redis/interface/database"
	"github.com/rong118/go_mini_redis/lib/logger"
	"github.com/rong118/go_mini_redis/lib/sync/atomic"
	"github.com/rong118/go_mini_redis/resp/connection"
	"github.com/rong118/go_mini_redis/resp/parser"
	"github.com/rong118/go_mini_redis/resp/reply"
)

type RespHandler struct {
  activeConn sync.Map
  closing atomic.Boolean
  db dataBaseFace.DataBase
}

func MakeHandler() *RespHandler {
  var db dataBaseFace.DataBase
  db = database.NewDataBase()
  return &RespHandler{
    db : db,
  }
}

func (r *RespHandler) closeClient(client *connection.Connection) {
  _ = client.Close()
  r.db.AfterClientClose(client)
  r.activeConn.Delete(client)
}

func (r *RespHandler) Handle(ctx context.Context, conn net.Conn){
  if r.closing.Get() {
    _ = conn.Close()
  }
  client := connection.NewConn(conn)
  r.activeConn.Store(client, struct{}{})

  ch := parser.ParserStream(conn)

  //监听管道
  for payload := range ch {
    if payload.Err != nil { // 异常处理
      // 用户关闭请求
      if payload.Err == io.EOF || payload.Err == io.ErrUnexpectedEOF || strings.Contains(payload.Err.Error(), "use of closed network connection") {
        r.closeClient(client)
        logger.Info("Connection closed: " + client.RemoteAddr().String())
        return
      }

      // protocol error
      errReply := reply.MakeErrReply(payload.Err.Error())
      err := client.Write(errReply.ToBytes())

      if err != nil {
        r.closeClient(client)
        logger.Info("connection closed " + client.RemoteAddr().String())
        return
      }
      continue
    }

    // exec
    if payload.Data == nil {
      continue
    }
    
    reply, ok := payload.Data.(*reply.MultiBulkReply)


    if !ok {
      var unkownErrReplyBytes = []byte("-ERR unknown\r\n")
      _ = client.Write(unkownErrReplyBytes)

      return
    }
    result := r.db.Exec(client, reply.Args)
    if result != nil {
      _ = client.Write(result.ToBytes())
    }else {
      var unkownErrReplyBytes = []byte("-ERR unknown\r\n")
      _ = client.Write(unkownErrReplyBytes)
    }
  }
}

func (r *RespHandler) Close() error {
  logger.Info("handler shutdown...")
  r.closing.Set(true)
  r.activeConn.Range(func(key interface{}, value interface{}) bool {
      client := key.(*connection.Connection)
      _ = client.Close()
      return true
    })

  r.db.Close()

  return nil
}
