package main 
import (
    "myhttp1.1/redisopt"
    "fmt"
    "log"
    "os"
)

func main(){
    method,_ := redisopt.CommandTable["GET"]
    get,ok := method.(redisopt.Get)
    if !ok{
        fmt.Println("get method wrong")
        return
    }
    //conn := redisopt.RedisDB.GetConn()
    //get.Init("GET",":6379")
    //argv := make([]interface{},0)
    //argv = append(argv,"foo")

    argv := "foo"
    res,err := get.Run(1,argv)
    if err != nil {
        fmt.Println(err)
        return
    }
    filename := "test.log"
    logFile,err := os.Create(filename)

    defer logFile.Close()

    debuglog := log.New(logFile,"[Debug]",log.LstdFlags)
    debuglog.Println("a devug message here")

    debuglog.SetPrefix("[info]")
    debuglog.Println("a debug message here")

    debuglog.SetFlags(debuglog.Flags() | log.LstdFlags)
    debuglog.Println("a debug message here")

    fmt.Println(res)
}
