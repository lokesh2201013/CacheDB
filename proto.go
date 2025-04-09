package main

import ("fmt"
"github.com/tidwall/resp"
"bytes"
"io"
"log"
)

const(
	CommandSET="SET"
)

type SetCommand struct {
	key, val string
}
type Command interface{}

func parseCommand(raw string) (Command, error) {
	rd:=resp.NewReader(bytes.NewBufferString(raw))
	for{
		v,_,err:=rd.ReadValue()

		if err==io.EOF{
			break
		}
		if err!=nil{
			log.Fatal(err)
		}

		fmt.Printf("Read %s\n",v.Type())

	

		if v.Type() == resp.Array{
			for _,value:= range v.Array(){
				switch value.String(){
				case CommandSET:
					fmt.Println(len(v.Array()))
					if len(v.Array())!=3{
						return nil, fmt.Errorf("Invalid command length")
					}
					cmd:=SetCommand{
						
						key: v.Array()[1].String(),
						val: v.Array()[2].String(),
					}
					return cmd, nil
				}
			}
	}
	
	return nil,fmt.Errorf("Invalid or unknown commmand recived")

}
   return nil,fmt.Errorf("Invalid or unknown commmand recived")
}