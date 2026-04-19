package server

import (
	"log"
	"syscall"

	"github.com/adityjoshi/iDB/config"
)

var connected_clients int = 0

func AsyncTcpServer() error {
	log.Println("Asysn Tcp Started on", config.Host, config.Port)

	var max_clients int = 20000

	var events []syscall.Kevent_t = make([]syscall.Kevent_t, max_clients)
}
