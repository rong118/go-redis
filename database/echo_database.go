package database

import (
	"go_mini_redis/interface/resp"
	"go_mini_redis/resp/reply"
)

type EchoDatabase struct {
}

func NewEchoDataBase() *EchoDatabase {
	return &EchoDatabase{}
}

// echo reply for mocking
func (e *EchoDatabase) Exec(client resp.Connection, arg [][]byte) resp.Reply {
	return reply.MakeMultiBulkReply(arg)
}

func (e *EchoDatabase) Close() {

}

func (e *EchoDatabase) AfterClientClose(c resp.Connection) {

}
