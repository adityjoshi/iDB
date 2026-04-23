package server

import (
	"log"
	"net"
	"syscall"
	"time"

	"github.com/adityjoshi/iDB/config"
	"github.com/adityjoshi/iDB/core"
)

var connectedClients int = 0
var cronFrequency time.Duration = 1 * time.Second
var lastCronExecTime time.Time = time.Now()

func AsyncTcpServer() error {
	log.Println("Async TCP Started on", config.Host, config.Port)

	maxClients := 20000
	events := make([]syscall.Kevent_t, maxClients)

	// Create socket
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(serverFD)

	// Set non-blocking
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		return err
	}

	// Bind
	ip4 := net.ParseIP(config.Host).To4()
	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		return err
	}

	// Listen
	if err = syscall.Listen(serverFD, syscall.SOMAXCONN); err != nil {
		return err
	}

	// Create kqueue
	kqFD, err := syscall.Kqueue()
	if err != nil {
		return err
	}
	defer syscall.Close(kqFD)

	// Register server FD
	serverEvent := syscall.Kevent_t{
		Ident:  uint64(serverFD),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD,
	}

	if _, err = syscall.Kevent(kqFD, []syscall.Kevent_t{serverEvent}, nil, nil); err != nil {
		return err
	}

	for {

		// cron events

		if time.Now().After(lastCronExecTime.Add(cronFrequency)) {
			core.DeleteExpiredKey()
			lastCronExecTime = time.Now()
		}

		nEvents, err := syscall.Kevent(kqFD, nil, events, nil)
		if err != nil {
			return err
		}

		for i := 0; i < nEvents; i++ {
			fd := int(events[i].Ident)

			// 🔹 New connection
			if fd == serverFD {
				clientFD, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println("accept error:", err)
					continue
				}

				connectedClients++

				// Set client non-blocking
				syscall.SetNonblock(clientFD, true)

				// Register client FD
				clientEvent := syscall.Kevent_t{
					Ident:  uint64(clientFD),
					Filter: syscall.EVFILT_READ,
					Flags:  syscall.EV_ADD,
				}

				if _, err := syscall.Kevent(kqFD, []syscall.Kevent_t{clientEvent}, nil, nil); err != nil {
					log.Println("kqueue register error:", err)
					syscall.Close(clientFD)
					continue
				}

			} else {
				// Handle client request
				comm := core.FileDescriptor{FD: fd}

				cmd, err := readCommand(comm)

				if err != nil {
					syscall.Close(fd)
					connectedClients--
					continue
				}

				respond(cmd, comm)
			}
		}
	}
}
