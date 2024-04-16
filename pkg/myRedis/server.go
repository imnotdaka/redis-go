package pkg

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	listener net.Listener
	protocol string
	ip       string
	cache    Storage
}

func NewServer(protocol string, ip string) Server {
	cache := newStorage()
	return Server{protocol: protocol, ip: ip, cache: cache}
}

func (s *Server) Start() error {
	listener, err := net.Listen(s.protocol, s.ip)
	if err != nil {
		fmt.Print("Failed to bind to port 6379")
		return err
	}

	s.listener = listener
	s.Listen()
	return nil
}

func (s Server) Listen() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Print("Error accepting connection: ", err.Error())
			return err
		}
		go s.HandleClient(conn)
	}
}

func (s Server) HandleClient(conn net.Conn) error {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			return err
		}
		commands := protocolParser(buf)
		switch strings.ToUpper(commands[0]) {
		case "ECHO":
			_, err = conn.Write([]byte(createBulkResponse(commands[1])))
			if err != nil {
				return err
			}
		case "SET":
			duration := 10000
			if len(commands) > 3 {
				duration, err = strconv.Atoi(commands[4])
				if err != nil {
					conn.Write([]byte(createBulkResponse("error")))
					return nil
				}
			}
			s.cache.SetValue(commands[1], commands[2], duration)
			_, err = conn.Write([]byte(createBulkResponse("OK")))
			if err != nil {
				return err
			}
		case "GET":
			s := s.cache.GetValue(commands[1])
			if s == "" {
				_, err = conn.Write([]byte(createNullResponse()))
				if err != nil {
					return err
				}
				return nil
			}
			_, err = conn.Write([]byte(createBulkResponse(s)))
			if err != nil {
				return err
			}
		case "PING":
			_, err = conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				return err
			}
		}

	}
}
