package core

import (
	"errors"
	"net"
)


func evalPING(args []string, conn net.Conn)(error){
    var b []byte

    if len(args)>=2{
        return errors.New("ERR wrong no. of args")
    }

    if len(args)==0{
        b=Encode("PONG",true)
    }else{
        b=Encode(args[0],false)
    }

    _,err:=conn.Write(b)
    return err
}

func EvalAndRespond(cmd *BredisCmd, conn net.Conn) error{
    switch cmd.Cmd{
    case "PING":    
        return evalPING(cmd.Args,conn)
    default:
        return evalPING(cmd.Args,conn)
    }
}