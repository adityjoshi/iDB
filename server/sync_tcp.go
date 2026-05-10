package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/adityjoshi/iDB/config"
	"github.com/adityjoshi/iDB/core"
)

func toArrayString(ai []interface{}) ([]string, error) {
	as := make([]string, len(ai))

	for i := range as {
		as[i] = ai[i].(string)
	}
	return as, nil
}

func readCommand(c io.ReadWriter) (core.RedisCmds, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return nil, err
	}

	values, err := core.Decode(buf[:n])
	if err != nil {
		return nil, err
	}
	var cmds []*core.RedisCmd = make([]*core.RedisCmd, 0)

	for _, value := range values {
		tokens, err := toArrayString(value.([]interface{}))
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, &core.RedisCmd{
			Cmd:  strings.ToUpper(tokens[0]),
			Args: tokens[1:],
		})
	}
	return cmds, nil
}

func respondError(err error, c io.ReadWriter) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(cmds core.RedisCmds, c io.ReadWriter) {

	core.EvaluateAndResponse(cmds, c)
}

func RunTcpServer() {
	log.Println("Synchronous TCP server started on", config.Host, config.Port)

	var connected_clients int = 0

	listner, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		log.Println("error", err)
	}

	for {
		/*
		*blocking call will wait for the clients to connect
		* */
		c, err := listner.Accept()
		if err != nil {
			log.Println("error", err)
		}

		connected_clients += 1

		log.Println("client connected to the server with address:", c.RemoteAddr(), "Concurrent clients -> ", connected_clients)

		for {

			cmds, err := readCommand(c)
			if err != nil {
				c.Close()
				connected_clients -= 1

				if err == io.EOF {
					break
				}
			}

			respond(cmds, c)

		}
	}

}
