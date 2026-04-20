package server

import (
	"log"
	"net"
	"syscall"

	"github.com/adityjoshi/iDB/config"
	"github.com/adityjoshi/iDB/core"
	"github.com/adityjoshi/iDB/server"
	"github.com/twitchyliquid64/golang-asm/sys"
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

	kqFD, err := syscall.Kqueue()

	if err != nil {
		return err
	}

	defer syscall.Close(kqFD)

	var socketServerEvents syscall.Kevent_t = syscall.Kevent_t{
		Ident:  uint64(serverFD),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD,
	}

	if _, err = syscall.Kevent(kqFD, []syscall.Kevent_t{socketServerEvents}, nil, nil); err != nil {
		return err
	}

	for {

		newEvents, err := syscall.Kevent(kqFD, nil, events[:], nil)
		if err != nil {
			return err
		}

		for i := 0; i < newEvents; i++ {
			if int(events[i].Ident) == serverFD {
				fd, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println("error", err)
					continue
				}

				connected_clients++
				syscall.SetNonblock(serverFD, true)

				var socketClientEvents syscall.Kevent_t = syscall.Kevent_t{
					Ident:  uint64(serverFD),
					Filter: syscall.EVFILT_READ,
					Flags:  syscall.EV_ADD,
				}

				if _, err := syscall.Kevent(fd, []syscall.Kevent_t{socketClientEvents}, nil, nil); err != nil {
					log.Fatal(err)
				} else {
					comm := core.FileDescriptor{FD: int(events[i].FD)}
					cmd, err := readCommand(comm)
					if err != nil {
						syscall.Close(int(events[i].FD))
						connected_clients -= 1
						continue
					}
					respond(cmd, comm)
				}
			}
		}
	}

}
