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

func readCommand(c io.ReadWriter) (*core.RedisCmd, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return nil, err
	}

	tokens, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}

	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil

}

func respondError(err error, c io.ReadWriter) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(cmd *core.RedisCmd, c io.ReadWriter) {

	err := core.EvaluateAndResponse(cmd, c)
	if err != nil {
		respondError(err, c)
	}
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

			cmd, err := readCommand(c)
			if err != nil {
				c.Close()
				connected_clients -= 1

				if err == io.EOF {
					break
				}
			}

			respond(cmd, c)

		}
	}

}
