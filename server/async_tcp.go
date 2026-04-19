package server

import (
	"log"
	"net"
	"syscall"

	"github.com/adityjoshi/iDB/config"
)

var connected_clients int = 0

func AsyncTcpServer() error {
	log.Println("Asysn Tcp Started on", config.Host, config.Port)

	var max_clients int = 20000

	var events []syscall.Kevent_t = make([]syscall.Kevent_t, max_clients)

	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}

	defer syscall.Close(serverFD)

	if err = syscall.SetNonblock(serverFD, true); err != nil {
		return err
	}

	ip4 := net.ParseIP(config.Host)

	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		return err
	}

	/*
	* Async event creation
	* */

}
