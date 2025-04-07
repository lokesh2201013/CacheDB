package main
 
import ("fmt"
"net"
"log/slog"
)

type Config struct{
	ListenAddr string
	}

type Server struct{
	Config
	peers map[*Peer]bool
	ln net.Listener
	addPeerCh chan *Peer
	}
	
func NewServer(cfg Config) *Server{

	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = ":3000"
	}

	return &Server{
		Config: cfg,
		peers: make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
	}
}

func (s *Server) Start() error{
	ln,err:=net.Listen("tcp",s.ListenAddr)
	if err != nil {
		return err
	}

	s.ln = ln

	go s.loop()
    
    return  s.acceptLoop()

}

func (s *Server) loop() {
	for{
		select{
			case peer := <-s.addPeerCh:
				s.peers[peer] = true
			default:
				fmt.Println("default")	
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

}



func main() {
	fmt.Println("Hello, World!")
}
