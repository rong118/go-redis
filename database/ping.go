package database

import (
	"github.com/rong118/go_mini_redis/interface/resp"
	"github.com/rong118/go_mini_redis/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}

func init() {
	RegisterCommand("ping", Ping, 1)
}
