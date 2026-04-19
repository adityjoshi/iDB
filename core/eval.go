package core

import (
	"io"
)

func evalPing(args []string, c io.ReadWriter) error {
	var b []byte

	if len(args) >= 2 {
		c.Write([]byte("-Err wrong number of argumnents for ping command\r\n"))
	}

	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}

	_, err := c.Write(b)
	return err

}

func EvaluateAndResponse(cmd *RedisCmd, c io.ReadWriter) error {

	switch cmd.Cmd {

	case "PING":
		return evalPing(cmd.Args, c)

	default:
		return evalPing(cmd.Args, c)
	}

}
