package core

import (
    "log"
    "time"
)

func expireSample() float32{

    var limit int =20

    var expiredCount int =0

    for key,obj:=range store{
        if obj.ExpiresAt!=-1{
            limit--
            if obj.ExpiresAt<=time.Now().UnixMilli(){
                expiredCount++
                Del(key)
            }
            if limit==0{
                break
            }
        }
    }
    return float32(expiredCount)/float32(20.0)
}


func DeleteExpiredKeys(){
    for{
        frac:=expireSample()

        if frac<0.25{
            break
        }

        log.Println("deleting expired keys: total left keys: ", len(store))
    }
}