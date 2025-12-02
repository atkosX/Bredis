package server

import (
	"io"
	"log"
	"main/config"
	"net"
	"strconv"
)

func readCmnd(conn net.Conn) (string, error) {

	var buf []byte = make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil

}

func respond(cmnd string, conn net.Conn) error {
	if _, err := conn.Write([]byte(cmnd)); err != nil {
		return err
	}
	return nil

}

func RunSyncTCPServer() {
	log.Println("starting a sync TCP server on", config.Host, ":", config.Port)

	var conn_clients int = 0

	listener, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error accepting connection: ", err)
			continue
		}
		conn_clients += 1
		log.Println("new connection accepted: ", conn.RemoteAddr(), "total connections: ", conn_clients)

		for {
			cmd, err := readCmnd(conn)
			if err != nil {
				conn_clients -= 1
				conn.Close()
				log.Println("error reading command: ", err)
				break
			}
			if err == io.EOF {
				conn_clients -= 1
				conn.Close()
				log.Println("connection closed by client: ", conn.RemoteAddr())
				break
			}

			log.Println("command received: ", cmd)
			if err = respond(cmd, conn); err != nil {
				log.Println("error responding to command: ", err)
				break
			}
		}
	}

}
