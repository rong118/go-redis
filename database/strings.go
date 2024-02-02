package database

import (
	"github.com/rong118/go_mini_redis/interface/database"
	"github.com/rong118/go_mini_redis/interface/resp"
	"github.com/rong118/go_mini_redis/resp/reply"
)

// GET k1
func execGet(db *DB, args [][]byte) resp.Reply {
  key := string(args[0])
	entity, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
  bytes := entity.Data.([]byte)
  return reply.MakeBulkReply(bytes)
}

// SET k, v
func execSet(db *DB, args [][]byte) resp.Reply {
  key := string(args[0])
  val := args[1]

  entity := &database.DataEntity{
    Data: val,
  }

  db.PutEntity(key, entity)

  return reply.MakeOkReply()

  
}

// SETNX
func execSetnx(db *DB, args [][]byte) resp.Reply {
  key := string(args[0])
  val := args[1]

  entity := &database.DataEntity{
    Data: val,
  }

  result := db.PutIfAbsent(key, entity)

  return reply.MakeIntReply(int64(result))
}

// GETSET
func execGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := args[1]

  entity, ok := db.GetEntity(key)
  db.PutEntity(key, &database.DataEntity{Data: value})

  if !ok {
    return reply.MakeNullBulkReply()
  }

	return reply.MakeBulkReply(entity.Data.([]byte))
}

// STRLEN  k1 v => len
func execStrLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
  entity, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}

  bytes := entity.Data.([]byte)

	return reply.MakeIntReply(int64(len(bytes)))
}

// 
func init() {
  RegisterCommand("GET", execGet, 3)
  RegisterCommand("SET", execSet, 3)
  RegisterCommand("SETNX", execSetnx, 3)
  RegisterCommand("GETSET", execGetSet, 3)
  RegisterCommand("STRLEN", execStrLen, 3)
}
