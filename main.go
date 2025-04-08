package main

import (
	"fmt"
	"log/slog"
	"net"
)

type Config struct{
	ListenAddr string
	}

type Server struct{
	Config
	peers map[*Peer]bool
	ln net.Listener
	addPeerCh chan *Peer
	quitCh chan struct{}
	msgCh chan []byte
	}
	
func NewServer(cfg Config) *Server{

	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = ":3000"
	}

	return &Server{
		Config: cfg,
		peers: make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh: make(chan struct{}),
		msgCh: make(chan []byte),
	}
}

func (s *Server) Start() error{
	ln,err:=net.Listen("tcp",s.ListenAddr)
	if err != nil {
		return err
	}

	s.ln = ln
	slog.Info("server started", "listenAddr", s.ListenAddr)
	go s.loop()
    
    return  s.acceptLoop()

}

func (s *Server) loop() {
	for{
		select{
			case rawMsg := <-s.msgCh:

				fmt.Println("msg:", (rawMsg))
			case <-s.quitCh:
				return
			case peer := <-s.addPeerCh:
				s.peers[peer] = true

		}
	}
}

func (s *Server)acceptLoop() error{
   for {
	 conn,err:=s.ln.Accept()
	 if err != nil {
		slog.Error("accept error",err)
		continue
	 }
	 go s.handleConn(conn)
   }
}

func (s *Server)handleConn(conn net.Conn){
peer:= NewPeer(conn,s.msgCh)
 s.addPeerCh <- peer
slog.Info("new peer","remoteAddr",conn.RemoteAddr(),)
 //go peer.readLoop()

 if err := peer.readLoop(); err != nil {
		slog.Error("peer read loop error",err,"remoteAddr",conn.RemoteAddr())
	}
}



func main() {
	server:=NewServer(Config{})
	err:=server.Start()
	if err != nil {
		slog.Error("start server error",err)
	}
}
