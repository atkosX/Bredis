package server

import (
	"log"
	"main/config"
	"main/core"
	"net"
	"syscall"
)

func RunAsyncTCPServer() error {
	log.Println("starting async TCP server on", config.Host, config.Port)

	max_clients := 20000

	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_clients)

	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal("error creating socket: ", err)
		return err
	}

	defer syscall.Close(serverFD)

	if err = syscall.SetNonblock(serverFD, true); err != nil {
		log.Fatal("error setting nonblock: ", err)
		return err
	}

	ip4 := net.ParseIP(config.Host)
	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		log.Fatal("error binding socket: ", err)
		return err
	}

	if err = syscall.Listen(serverFD, max_clients); err != nil {
		log.Fatal("error listening socket: ", err)
		return err
	}

	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal("error creating epoll: ", err)
	}

	defer syscall.Close(epollFD)

	var socketServerEvent syscall.EpollEvent = syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFD),
	}

	if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, serverFD, &socketServerEvent); err != nil {
		log.Fatal("error adding socket to epoll: ", err)
		return err
	}

	con_clients := 0

	for {
		nevents, e := syscall.EpollWait(epollFD, events[:], -1)
		if e != nil {
			log.Fatal("error waiting for events: ", e)
			return e
		}

		for i := 0; i < nevents; i++ {
			if events[i].Fd == int32(serverFD) {
				fd, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Fatal("error accepting connection: ", err)
					continue
				}
				log.Println("added a new connection: ", fd)
				con_clients++
				syscall.SetNonblock(fd, true)

				var socketClientEvent syscall.EpollEvent = syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(fd),
				}

				if err := syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, fd, &socketClientEvent); err != nil {
					log.Fatal("error adding socket to epoll: ", err)
				}

			} else {
				comm := core.FDcommand{FD: int(events[i].Fd)}
				cmd, err := readCmd(comm)
				if err != nil {
					syscall.Close(int(events[i].Fd))
					log.Fatal("error reading command: ", err)
					con_clients--
					continue
				}
				respond(cmd, comm)
			}
		}
	}

}
