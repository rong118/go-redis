package database

import (
	"go_mini_redis/interface/resp"
)

type Cmdline = [][]byte

type DataBase interface {
	Exec(client resp.Connection, arg [][]byte) resp.Reply
	Close()
	AfterClientClose(c resp.Connection)
}

type DataEntity struct {
	Data interface{}
}
