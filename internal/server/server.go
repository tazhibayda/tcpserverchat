package server

import (
	"context"
	"log"
	"net"
	"sync"

	"tcpserverchat/internal/client"
	"tcpserverchat/internal/dispatcher"
)

type Server struct {
	listener net.Listener
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func New(addr string) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	return &Server{listener: l, ctx: ctx, cancel: cancel}
}

func (s *Server) Start() error {
	d := dispatcher.New(s.ctx)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		d.Run()
	}()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.ctx.Done():
				return nil
			default:
				log.Println("accept:", err)
				continue
			}
		}
		s.wg.Add(1)
		go func(c net.Conn) {
			defer s.wg.Done()
			client.Handle(c, d)
		}(conn)
	}
}

func (s *Server) Stop() {
	s.cancel()
	s.listener.Close()
	s.wg.Wait()
}
