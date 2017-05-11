package redisopt

import (
//	"fmt"
	"log"
	"os"
)
const (
    client = iota
    server
)
var Logmapser = make(map[string]*log.Logger)
var Logmapcli = make(map[string]*log.Logger)
func init(){
    InitLog("logserver",server)
    InitLog("logcli",client)
}
func InitLog(filename string,rule int) {
	//初始化日志信息
	file, err1 := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err1 != nil {
		log.Fatal(err1)
	}
	//defer file.Close()
	errlog := log.New(file, "[err]", log.Lshortfile|log.LstdFlags)
	infolog := log.New(file, "[message]", log.Lshortfile|log.LstdFlags)

    if rule == client{
        Logmapcli["errlog"] = errlog
        Logmapcli["infolog"] = infolog
    }else if rule == server{
        Logmapser["errlog"] = errlog
        Logmapser["infolog"] = infolog
    }
}
