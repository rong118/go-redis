package reply

/* ---- Reply Consts---- */

// PongReply
type PongReply struct {}

var pongBytes = []byte("+PONG\r\n")

func (r* PongReply) ToBytes() []byte {
  return pongBytes
}

func MakePongReply() *PongReply {
  return &PongReply{}
}

// OkReply
type OkReply struct{}

var okBytes = []byte("+OK\r\n")

func (r* OkReply) ToBytes() []byte {
  return okBytes
}

func MakeOkReply() *OkReply {
  return &OkReply{}
}

// Empty Bulk Reply
type NullBulkReply struct{}

var nullBulkBytes = []byte("$-1\r\n")

func (r* NullBulkReply) ToBytes() []byte {
  return nullBulkBytes
}

func MakeNullBulkReply() *NullBulkReply {
  return &NullBulkReply{}
}

// Empty Multi-Bulk Reply
type EmptyMultiBulkReply struct{}

var emptyMultiBulkBytes = []byte("*0\r\n")

func (r* EmptyMultiBulkReply) ToBytes() []byte {
  return emptyMultiBulkBytes
}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
  return &EmptyMultiBulkReply{}
}

// No reply
type NoReply struct {}

var noReplyBytes = []byte("")

func (r* NoReply) ToBytes() []byte {
  return noReplyBytes
}

func MakeNoReply() *NoReply {
  return &NoReply{}
}

