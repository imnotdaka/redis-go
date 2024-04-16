package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	LF    = "\r\n"
	array = "*"
)

type Server struct {
	listener net.Listener
	protocol string
	ip       string
	cache    Storage
}

func NewServer(protocol string, ip string, cache Storage) Server {
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
		switch commands[0] {
		case "echo":
			_, err = conn.Write([]byte(createBulkResponse(commands[1])))
			if err != nil {
				return err
			}
		case "set":
			duration := 100
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
		case "get":
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
		case "ping":
			_, err = conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				return err
			}
		}

	}
}

func protocolParser(input []byte) []string {
	inputStr := string(input)
	var commands []string
	if strings.HasPrefix(inputStr, array) {
		commands = ParseArray(inputStr)
	}

	return commands
}
func ParseArray(inputStr string) []string {
	var res []string
	end := strings.Index(inputStr, "\r\n")
	s := inputStr[end+2:]
	for i := 0; i < len(s); i++ {
		if string(s[i]) == "$" {
			j := strings.Index(s[i:], "\r\n")
			n, _ := strconv.Atoi(s[i+1 : j+i])
			word := s[j+i+2 : j+i+2+n]
			res = append(res, word)
		}
		// $4\r\necho\r\n$3\r\nhey\r\n
	}
	return res
}

func createBulkResponse(s string) string {

	return "$" + strconv.Itoa(len(s)) + LF + s + LF
}

func createNullResponse() string {

	return "$-1" + LF
}
