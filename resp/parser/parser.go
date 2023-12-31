package parser

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/rong118/go_mini_redis/interface/resp"
	"github.com/rong118/go_mini_redis/lib/logger"
	"github.com/rong118/go_mini_redis/resp/reply"
)

type Payload struct {
  Data resp.Reply
  Err error
}

type readState struct {
  readingMultiLine bool
  expectedArgsCount int
  msgType byte
  args [][]byte
  bulkLen int64
}

func (r *readState) finished() bool {
  return r.expectedArgsCount > 0 && len(r.args) == r.expectedArgsCount
}

func ParserStream(reader io.Reader) <-chan *Payload {
  ch := make(chan *Payload)
  go _parser(reader, ch)
  return ch
}

func _parser(reader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
  bufReader := bufio.NewReader(reader)
  var state readState
  var err error
  var msg []byte
	for true {
    var ioErr bool
    msg, ioErr, err = readLine(bufReader, &state)
    if err != nil {
      if ioErr {
        	ch <- &Payload{
            Err: err,
          }
  		  	close(ch)
	    		return
      }else {
        ch<-&Payload{
          Err: err,
        }
        state = readState{}
        continue
      }
    }

    if !state.readingMultiLine {
      if msg[0] == '*' { //*3/r/n
        err := parseMultiBulkHeader(msg, &state)
        if err != nil {
	        ch <- &Payload{
            Err: errors.New("protocol error 1: " + string(msg)),
          }
          state = readState{}
          continue
        }

        if state.expectedAtgsCount == 0 {
          ch <- &Payload{
            Data: &reply.EmptyMultiBulkReply{},
          }
          state = readState{}
          continue
        }
      } else if msg[0] == '$' { // $3\r\nSET\r\n
        err := parseBulkHeader(msg, &state)
        if err != nil {
	        ch <- &Payload{
            Err: errors.New("protocol error 2: " + string(msg)),
          }
          state = readState{}
          continue
        }

        if state.bulkLen == -1 {  // -1/r/n
          ch <- &Payload{
            Data: &reply.EmptyMultiBulkReply{},
          }
          state = readState{}
          continue
        }
      } else {
        result, err := parserSingleLineReply(msg)
        ch <- &Payload{
            Data: result,
            Err: err,
          }
          state = readState{}
          continue
      }
    } else {
      err := readBody(msg, &state)
      if err != nil {
        ch <- &Payload{
          Err: errors.New("protocol error 3: " + string(msg)),
        }
        state = readState{}
        continue
      }

      if state.finished() {
        var result resp.Reply
        if state.msgType == '*' {
          result = reply.MakeMultiBulkReply(state.args)
        } else if state.msgType == '$' {

          result = reply.MakeBulkReply(state.args[0])
        }
        ch <- &Payload{
          Data: result,
          Err: err,
        }
        state = readState{}
        break
      }
    }
	}
}


func readLine(bufReader *bufio.Reader, state *readState)([]byte, bool, error) {
  var msg []byte
  var err error
  if state.bulkLen == 0 {  // 1.没有读取$，以\r\n为准 
    msg, err = bufReader.ReadBytes('\n')
    
    if err != nil {
      return nil, true, err
    }
    if len(msg) == 0 || msg[len(msg) - 2] != '\r' {
      return nil, false, errors.New("protocol error: 4" + string(msg))
    }
  }else{ // 2. 读取了$， 以字符个数为准
    msg = make([]byte, state.bulkLen + 2) //len + \r\n

    _, err := io.ReadFull(bufReader, msg)
    if err != nil {
      return nil, true, err
    }
    if len(msg) == 0 || msg[len(msg) - 2] != '\r' ||  msg[len(msg) - 1] != '\n' {
      return nil, false, errors.New("protocol error: 5 " + string(msg))
    } 
    state.bulkLen = 0
  }
  return msg, false, nil
}

// *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
func parseMultiBulkHeader(msg []byte, state *readState) error {
  var err error
  var expectedLine uint64
  expectedLine, err = strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 32)
  if err != nil {
    return  errors.New("protocol error: 6" + string(msg))
  }

  if expectedLine == 0 {
    state.expectedAtgsCount = 0
    return nil
  }else if expectedLine > 0{
    state.msgType = msg[0]
    state.readingMultiLine = true
    state.expectedAtgsCount = int(expectedLine)
    state.args = make([][]byte, 0, expectedLine)
    return nil
  }else{
    return errors.New("protocol error: 7" + string(msg))
  }
}

// $4\r\nPING\r\n
func parseBulkHeader(msg []byte, state *readState) error {
  var err error
  state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 32)
  if err != nil {
    return  errors.New("protocol error: 8" + string(msg))
  }

  if state.bulkLen == -1 {
    return nil
  }else if state.bulkLen > 0{
    state.msgType = msg[0]
    state.readingMultiLine = true
    state.expectedAtgsCount = 1
    state.args = make([][]byte, 0, 1)
    return nil
  }else{
    return errors.New("protocol error: 9" + string(msg))
  }
}

// +OK\r\n -err\r\n :5\r\n
func parserSingleLineReply(msg []byte)(resp.Reply, error) {
  str := strings.TrimSuffix(string(msg), "\r\n")
  var result resp.Reply
  switch msg[0] {
  case '+':
    result = reply.MakeStatusReply(str[1:])
  case '-':
    result = reply.MakeErrReply(str[1:])
  case ':':
    val, err := strconv.ParseInt(str[1:], 10, 64)
    if err != nil {
      return nil, errors.New("protocol error: 10" + string(msg))
    }
    result = reply.MakeIntReply(val)
  }
  return result, nil
}

// PING\r\n
func readBody(msg []byte, state *readState) error {
  line := msg[0: len(msg) - 2]
  var err error
  // $3
  if line[0] == '$'{
    state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
    if err != nil {
      return errors.New("protocol error: 11" + string(msg))
    }

    // $0\r\n
    if state.bulkLen <= 0 {
      state.args = append(state.args, []byte{})
      state.bulkLen = 0
    }
  } else {
    state.args = append(state.args, line)
  }

  return nil
}
