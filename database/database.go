package database

import (
	"strconv"
	"strings"

	"github.com/rong118/go_mini_redis/config"
	"github.com/rong118/go_mini_redis/interface/resp"
	"github.com/rong118/go_mini_redis/lib/logger"
	"github.com/rong118/go_mini_redis/resp/reply"
)

type Database struct {
	dbSet []*DB
}

func NewDataBase() *Database {
	database := &Database{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	database.dbSet = make([]*DB, config.Properties.Databases)

	for i := range database.dbSet {
		db := makeDB()
		db.index = i
	}

	return database
}

func (database *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) != 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(client, database, args[1:])
	}

	dbIndex := client.GetDBIndex()
	db := database.dbSet[dbIndex]
	return db.Exec(client, args)
}

func (e *Database) Close() {

}

func (e *Database) AfterClientClose(c resp.Connection) {

}

// select 1
func execSelect(c resp.Connection, database *Database, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR invalid DB index")
	}

	if dbIndex >= len(database.dbSet) {
		return reply.MakeErrReply("ERR DB index is out of range")
	}

	c.SelectDB(dbIndex)
	return reply.MakeOkReply()
}
