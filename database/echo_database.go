package database

import (
	"github.com/rong118/go_mini_redis/interface/resp"
	"github.com/rong118/go_mini_redis/resp/reply"
)

type EchoDatabase struct {

}

func NewEchoDataBase() *EchoDatabase {
  return &EchoDatabase{}
}

// echo reply for mocking 
func(e *EchoDatabase) Exec(client resp.Connection, arg []byte) resp.Reply {
  return reply.MakeBulkReply(arg)
}

func(e *EchoDatabase) Close() {
  
}

func(e *EchoDatabase) AfterClientClose(c resp.Connection) {

}
