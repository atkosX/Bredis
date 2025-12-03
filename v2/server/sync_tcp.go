package server

import (
	"fmt"
	"io"
	"log"
	"main/config"
	"main/core"
	"net"
	"strconv"
	"strings"
)

func readCmd(conn net.Conn) (*core.BredisCmd, error) {

	var buf []byte = make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	tokens,err:=core.DecodeArrayString(buf[:n])
	if err!=nil{
		return nil,err
	}

	return &core.BredisCmd{
		Cmd:strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	},nil

}

func respondError(err error, conn net.Conn){
	conn.Write([]byte(fmt.Sprintf("-%s\r\n",err)))
}

func respond(cmd *core.BredisCmd, conn net.Conn) error {
	err:=core.EvalAndRespond(cmd,conn)
	
	if err!=nil{
		respondError(err, conn)
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
			cmd, err := readCmd(conn)
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
