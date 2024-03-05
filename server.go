package main

import (
	"log"
	"net"
)

type Server struct {
	*Aof
	l        net.Listener
	handlers map[string]func([]Value)
}

func NewServer(dbPath string, port string, handlers map[string]func([]Value)) *Server {
	aof, err := NewAof(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	return &Server{aof, listener, handlers}
}

func (s *Server) Start() {
	//read all cmds from aof file and execute them
}

func (s *Server) handleConnection(c net.Conn) {
	defer c.Close()


}
