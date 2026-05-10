package core

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"time"
)

var RESP_NIL []byte = []byte("$-1\r\n")
var RESP_OK []byte = []byte("+OK\r\n")
var RESP_ZERO []byte = []byte(":0\r\n")
var RESP_ONE []byte = []byte(":1\r\n")
var RESP_MINUS_1 []byte = []byte(":-1\r\n")
var RESP_MINUS_2 []byte = []byte(":-2\r\n")

func evalPing(args []string) []byte {
	var b []byte

	if len(args) >= 2 {
		return Encode(errors.New("Err wrong number of argumnents for ping command"), false)
	}

	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}

	return b

}

func evalSet(args []string) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("(error) invalid number of arguments"), false)
	}

	var key, val string
	var expiraryDuration int64 = -1

	key, val = args[0], args[1]

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			i++
			if i == len(args) {
				return Encode(errors.New("(error) syntax error"), false)
			}
			expiraryDurationMS, err := strconv.ParseInt(args[3], 10, 64)

			if err != nil {
				return Encode(errors.New("(error) this value is not an integer or out of range"), false)
			}
			expiraryDuration = expiraryDurationMS * 1000
		default:
			return Encode(errors.New("(error) syntax error"), false)
		}
	}

	Put(key, NewObj(val, expiraryDuration))
	return RESP_OK
}

func evalGet(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) wrong number of arguments for the get command"), false)
	}

	var key string = args[0]

	Object := Get(key)

	if Object == nil {
		return RESP_NIL
	}

	if Object.ExpiresAt != -1 && Object.ExpiresAt <= time.Now().UnixMilli() {
		return RESP_NIL
	}

	return Encode(Object.value, false)
}

func evalTTL(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) Wrong number of arguments for the TTL Command"), false)
	}

	var key string = args[0]

	Object := Get(key)

	// returns -2 if the key doesn't exist
	if Object == nil {
		return RESP_MINUS_2
	}
	// return -1 if the key exist but has no associated expire
	if Object.ExpiresAt == -1 {
		return RESP_MINUS_1
	}

	/*
	* Compute the time remaining for the key
	* */
	durationMS := Object.ExpiresAt - time.Now().UnixMilli()

	if durationMS < 0 {
		return RESP_MINUS_2
	}

	return Encode(int64(durationMS/1000), false)

}

func evalDel(args []string) []byte {
	var countDelete int = 0

	for _, key := range args {
		if ok := Del(key); ok {
			countDelete++
		}
	}
	return Encode(countDelete, false)
}

func evalExpire(args []string) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("(error) invalid number of arguments for expire command"), false)
	}

	var key string = args[0]
	expiraryDurationSec, err := strconv.ParseInt(args[1], 10, 64)

	if err != nil {
		return Encode(errors.New("(error) Invalid value either it is not integer or out of range"), false)
	}

	Object := Get(key)

	if Object == nil {
		return RESP_ZERO
	}

	Object.ExpiresAt = time.Now().UnixMilli() + expiraryDurationSec*1000

	// send 1 if the expiration is set
	return RESP_ONE
}

func evalBGREWRITEAOF(args []string) []byte {
	DUMPALLAOF()
	return RESP_OK
}

func EvaluateAndResponse(cmds RedisCmds, c io.ReadWriter) {

	var response []byte
	buf := bytes.NewBuffer(response)

	for _, cmd := range cmds {
		switch cmd.Cmd {

		case "PING":
			buf.Write(evalPing(cmd.Args))

		case "SET":
			buf.Write(evalSet(cmd.Args))

		case "GET":
			buf.Write(evalGet(cmd.Args))

		case "TTL":
			buf.Write(evalTTL(cmd.Args))

		case "DEL":
			buf.Write(evalDel(cmd.Args))

		case "BGREWRITEAOF":
			buf.Write(evalBGREWRITEAOF(cmd.Args))

		case "EXPIRE":
			buf.Write(evalExpire(cmd.Args))
		default:
			buf.Write(evalPing(cmd.Args))
		}
	}
	c.Write(buf.Bytes())
}
