package database

import (
	"strings"

	"github.com/rong118/go_mini_redis/datastruct/dict"
	"github.com/rong118/go_mini_redis/interface/database"
	"github.com/rong118/go_mini_redis/interface/resp"
	"github.com/rong118/go_mini_redis/resp/reply"
)

type DB struct {
	index int
	data  dict.Dict
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply

type CmdLine = [][]byte

func makeDB() *DB {
	db := &DB{
		data: dict.MakeSyncDict(),
	}
	return db
}

func (db *DB) Exec(c resp.Connection, cmdline CmdLine) resp.Reply {
	//PING SET ...
	cmdName := strings.ToLower(string(cmdline[0]))

	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeErrReply("ERR unknown command " + cmdName)
	}

	if !validateArity(cmd.arity, cmdline) {
		return reply.MakeArgNumErrReply(cmdName)
	}

	fun := cmd.exector

	// SET K V ==> K V
	return fun(db, cmdline[1:])
}

// SET K V ==> arity = 3
// EXISTS K1 K2 K3 ... ==> arity = -2
func validateArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if argNum >= 0 {
		return argNum == arity
	}

	return argNum >= -1*arity
}

func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}

	entity, _ := raw.(*database.DataEntity)

	return entity, true
}

func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, ok := db.data.Get(key)
		if ok {
			db.Remove(key)
			deleted++
		}
	}

	return deleted
}

func (db *DB) Flush() {
	db.data.Clear()
}
