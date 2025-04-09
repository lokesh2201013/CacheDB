package main

import (
	"fmt"
	//"io"
	//"log"
	"testing"
	//"github.com/tidwall/resp" 
	//"bytes"
)

func TestProtocol(t *testing.T){
	raw:="*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
     cmd,err:=parseCommand(raw)
	 if err!=nil{
	   t.Fatal(err)
	}
	fmt.Println(cmd)
}