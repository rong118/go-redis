package database

import (
	"go_mini_redis/interface/resp"
	"go_mini_redis/lib/wildcard"
	"go_mini_redis/resp/reply"
)

// DEL
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}

	deleted := db.Removes(keys...)

	return reply.MakeIntReply(int64(deleted))
}

// EXISTS K1 K2 K3 ...
func execExists(db *DB, args [][]byte) resp.Reply {
	result := int64(0)
	for _, arg := range args {
		key := string(arg)
		_, exists := db.GetEntity(key)
		if exists {
			result++
		}
	}
	return reply.MakeIntReply(result)
}

// FLUSHDB
func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	return reply.MakeOkReply()
}

// TYPE
// execType returns the type of entity, including: string, list, hash, set and zset
func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeStatusReply("none")
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
		//TODO: support other data type
	}

	return &reply.UnkownErrReply{}
}

// RENAME
func execRename(db *DB, args [][]byte) resp.Reply {
	if len(args) != 2 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'rename' command")
	}
	src := string(args[0])
	dest := string(args[1])

	entity, ok := db.GetEntity(src)
	if !ok {
		return reply.MakeErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)

	return reply.MakeOkReply()
}

// RENAMENX
func execRenameNx(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])

	_, ok := db.GetEntity(dest)
	if ok {
		return reply.MakeIntReply(0)
	}

	entity, ok := db.GetEntity(src)
	if !ok {
		return reply.MakeErrReply("no such key")
	}
	db.Removes(src, dest) // clean src and dest with their ttl
	db.PutEntity(dest, entity)

	return reply.MakeIntReply(1)
}

// execKeys returns all keys matching the given pattern
func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern, err := wildcard.CompilePattern(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR illegal wildcard")
	}
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		if !pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}

		return true
	})

	return reply.MakeMultiBulkReply(result)
}

func init() {
	RegisterCommand("del", execDel, -2)
	RegisterCommand("exist", execExists, -2)
	RegisterCommand("FLUSHDB", execFlushDB, -1)
	RegisterCommand("TYPE", execType, 2)
	RegisterCommand("RENAME", execRename, 3)
	RegisterCommand("RENAMENX", execRenameNx, 3)
	RegisterCommand("KEYS", execKeys, 2)
}
