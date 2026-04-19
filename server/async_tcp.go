package server

import (
	"log"

	"github.com/adityjoshi/iDB/config"
)

var connected_clients int = 0

func AsyncTcpServer() error {
	log.Println("Asysn Tcp Started on", config.Host, config.Port)
}
