package core

import (
	"errors"
	"io"
	"strconv"
	"time"
)

var RESP_NIL []byte = []byte("$-1\r\n")

func evalPing(args []string, c io.ReadWriter) error {
	var b []byte

	if len(args) >= 2 {
		return errors.New("Err wrong number of argumnents for ping command\r\n")
	}

	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}

	_, err := c.Write(b)
	return err

}

func evalSet(args []string, c io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("(error) invalid number of arguments")
	}

	var key, val string
	var expiraryDuration int64 = -1

	key, val = args[0], args[1]

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			i++
			if i == len(args) {
				return errors.New("(error) syntax error")
			}
			expiraryDurationMS, err := strconv.ParseInt(args[3], 10, 64)

			if err != nil {
				return errors.New("(error) this value is not an integer or out of range")
			}
			expiraryDuration = expiraryDurationMS * 1000
		default:
			return errors.New("(error) syntax error")
		}
	}

	Put(key, NewObj(val, expiraryDuration))
	c.Write([]byte("+OK\r\n"))
	return nil
}

func evalGet(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) wrong number of arguments for the get command")
	}

	var key string = args[0]

	Object := Get(key)

	if Object == nil {
		c.Write(RESP_NIL)
	}

	if Object.ExpiresAt != -1 && Object.ExpiresAt <= time.Now().UnixMilli() {
		c.Write(RESP_NIL)
	}

	c.Write(Encode(Object.value, false))
	return nil
}

func evalTTL(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) Wrong number of arguments for the TTL Command")
	}

	var key string = args[0]

	Object := Get(key)

	if Object == nil {
		c.Write([]byte(":2\r\n"))
		return nil
	}

}

func EvaluateAndResponse(cmd *RedisCmd, c io.ReadWriter) error {

	switch cmd.Cmd {

	case "PING":
		return evalPing(cmd.Args, c)

	case "SET":
		return evalSet(cmd.Args, c)

	case "GET":
		return evalGet(cmd.Args, c)

	case "TTl":
		return evalTTL(cmd.Args, c)

	default:
		return evalPing(cmd.Args, c)
	}

}
