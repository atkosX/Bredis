package core

import (
	"errors"
	"io"
	"strconv"
	"time"
)

var RESP_NIL []byte=[]byte("$-1\r\n")

func evalPING(args []string, conn io.ReadWriter)(error){
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

func evalSET(args []string, conn io.ReadWriter)(error){
    if len(args)<=1{
        return errors.New("(error) wrong no. of args for set command")
    }

    var key, value string
    var exDurationMs int64 =-1

    key,value=args[0],args[1]

    for i:=2;i<len(args);i++{
        switch args[i]{
        case "EX", "ex":
            i++
            if i==len(args){
                return errors.New("(error) syntax err")
            }

            exDurationMs,err:=strconv.ParseInt(args[i],10,64)
            if err!=nil{
                return errors.New("(error) invalid duration")
            } 

            exDurationMs*=1000
            default:
                return  errors.New("(error) syntax err")
        }
    }

    obj:=NewObj(value,exDurationMs)
    Put(key,obj)
    conn.Write(Encode("OK",true))
    return nil
}

func evalGET(args []string, conn io.ReadWriter) error{
    if len(args)!=1{
        return errors.New("(error) invalid no. of args for get command")
    }

    var key string=args[0]
    obj:=Get(key)

    if obj==nil {
        conn.Write(RESP_NIL)
        return nil
    }

    if obj.ExpiresAt!=-1 && obj.ExpiresAt<=time.Now().UnixMilli(){
        conn.Write(RESP_NIL)
        return nil
    }

    conn.Write(Encode(obj.Value,false))
    return nil
}

func evalTTL(args []string, conn io.ReadWriter) error{
    if len(args)!=1{
        return errors.New("(error) invalid no. of args for ttl command")
    }

    var obj =Get(args[0])

    if obj==nil{
        conn.Write([]byte(":-2\r\n"))
        return nil
    }
    if obj.ExpiresAt==-1{
        conn.Write([]byte(":-1\r\n"))
        return nil
    }

    duratioMs:=obj.ExpiresAt-time.Now().UnixMilli()

    if duratioMs<0{
    conn.Write([]byte(":-2\r\n"))
    }

    conn.Write(Encode(int64(duratioMs/1000),false))
    return nil
}


func evalDEL(args []string,conn io.ReadWriter) error{
    if len(args)!=1{
        return errors.New("(error) invalid no. of args for del command")
    }

    var countdel int=0

    for _, key:=range args{
        if ok:=Del(key);ok{
            countdel++
        }
    }

    conn.Write(Encode(int64(countdel),false))
    return nil
}

func evalEXPIRE(args []string,conn io.ReadWriter) error{
    if len(args)!=2{
        return errors.New("(error) invalid no. of args for expire command")
    }

    var key string=args[0]
    var exDurationMs int64

    exDurationMs,err:=strconv.ParseInt(args[1],10,64)
    if err!=nil{
        return errors.New("(error) invalid duration : not an integer or out of range")
    }
    obj:=Get(key)
    if obj==nil{
        conn.Write(Encode(int64(0),false))
        return nil
    }
    obj.ExpiresAt=time.Now().UnixMilli()+exDurationMs*1000

    conn.Write(Encode(int64(1),false))

    return nil
}

func EvalAndRespond(cmd *BredisCmd, conn io.ReadWriter) error{
    switch cmd.Cmd{
    case "PING":    
        return evalPING(cmd.Args,conn)
    case "SET":
        return evalSET(cmd.Args,conn)
    case "GET":
        return evalGET(cmd.Args,conn)
    case "TTL":
        return evalTTL(cmd.Args,conn)
    case "DEL":
        return evalDEL(cmd.Args,conn)
    case "EXPIRE":
        return evalEXPIRE(cmd.Args,conn)
    default:
        return evalPING(cmd.Args,conn)
    }
}