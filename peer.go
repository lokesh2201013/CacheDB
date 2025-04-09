package main

import (
	//"fmt"
	"bufio"
	"io"
	"log/slog"
	"net"
)

type Peer struct {
	conn net.Conn	
	msgCh chan []byte
}

func NewPeer(conn net.Conn,msgCh chan []byte) *Peer {
	return &Peer{
		conn: conn,
		msgCh: msgCh,
	}
}

/*func(p *Peer) readLoop()error{
	buf:= make([]byte, 1024)
	
	for{
      n,err:=(p.conn.Read(buf))

	  if err != nil {
		slog.Error("read error",err)
		return err
	  } 
	
	  
	  msgBuf:=make([]byte, n)
	  copy(msgBuf,buf[:n])


	  p.msgCh<-msgBuf
	}
}*/
func (p *Peer) readLoop2() error {
	reader := bufio.NewReader(p.conn)

	for {
		line, err := reader.ReadString('\n') // wait until Enter is pressed
		if err != nil {
			if err == io.EOF {
				slog.Info("peer disconnected")
				return nil
			}
			slog.Error("read error", "err", err)
			return err
		}

		p.msgCh <- []byte(line)
	}
}


