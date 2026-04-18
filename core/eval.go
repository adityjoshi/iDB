package core

import "net"

func evalPing(args *RedisCmd, c net.Conn) {

}

func EvaluateAndResponse(cmd *RedisCmd, c net.Conn) {

	switch cmd.Cmd {

	case "PING":
		return evalPing(cmd.Args, c)

	default:
		return evalPing(cmd.Args, c)
	}

}
