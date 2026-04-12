package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/adityjoshi/iDB/config"
)

func readCommand(c net.Conn) (string, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func respond(cmd string, c net.Conn) error {
	if _, err := c.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}

func RunTcpServer() {
	log.Println("Synchronous TCP server started on", config.Host, config.Port)

	var connected_clients int = 0

	listner, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	for {
		/*
		*blocking call will wait for the clients to connect
		* */
		c, err := listner.Accept()
		if err != nil {
			panic(err)
		}
		connected_clients += 1

		log.Println("client connected to the server with address:", c.RemoteAddr(), "Concurrent clients -> ", connected_clients)

		for {

			cmd, err := readCommand(c)
			if err != nil {
				c.Close()
				connected_clients -= 1
				log.Println("client disconnected", c.RemoteAddr(), "Concurrent client -> ", connected_clients)

				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}

			log.Println("command", cmd)

			if err = respond(cmd, c); err != nil {
				log.Println("err write:", err)
			}
		}
	}

}
