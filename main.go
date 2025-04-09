package main

import (
	"fmt"
	"log"
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
//Intialises a server instance sets the listenADDR
//and intitalises the  Strict Server
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

//Listens on the specified address and starts the server
// Starts Goroutine of loop and returns s.acceptLoop
func (s *Server) Start() error{
	ln,err:=net.Listen("tcp",s.ListenAddr)
	if err != nil {
		return err
	}
    
	//Sets Server ln val
	s.ln = ln
	slog.Info("server started", "listenAddr", s.ListenAddr)
	go s.loop()
    
    return  s.acceptLoop()

}

func (s *Server)handleRawMessage(rawMsg []byte) error {
	cmd, err := parseCommand(string(rawMsg))

	if err != nil {
		return err
	}
     
	switch v:=cmd.(type){
		case SetCommand:
			slog.Info("set command", "key", v.key, "val", v.val)
	}
	return nil

}

//Goroutine that takes messages from the msgCh channel and handles them
// by calling handleRawMessage it takes in messages
//Listens for new peers and shut down signal
func (s *Server) loop() {
	for{
		select{
			case rawMsg := <-s.msgCh:
             if err:=   s.handleRawMessage(rawMsg); err != nil {
				slog.Error("handle raw message error","err",err)
				continue
			 }
				fmt.Println("msg:", (rawMsg))

			case <-s.quitCh:
				return

			case peer := <-s.addPeerCh:
				s.peers[peer] = true

		}
	}
}

//Accepst conn of Start set from Server ln
// Run Indefinately
//Gorotine HandleConn sends conn
func (s *Server)acceptLoop() error{
   for {
	 conn,err:=s.ln.Accept()
	 if err != nil {
		slog.Error("accept error","err",err)
		continue
	 }
	 go s.handleConn(conn)
   }
}

//
func (s *Server)handleConn(conn net.Conn){
peer:= NewPeer(conn,s.msgCh)
 s.addPeerCh <- peer
slog.Info("new peer","remoteAddr",conn.RemoteAddr(),)
 //go peer.readLoop()

 if err := peer.readLoop2(); err != nil {
		slog.Error("peer read loop error","err",err,"remoteAddr",conn.RemoteAddr())
	}
}



func main() {

	//Create neew server instance and starts it
	//cfg:=Config{ListenAddr:":3000"}
	go func(){
		server:=NewServer(Config{})
	     log.Fatal(server.Start())
	}()
	
	select {}
}
