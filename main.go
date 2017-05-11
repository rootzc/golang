package main
import (
    "fmt"
    "net/http"
    "strings"
    "git.apache.org/thrift.git/lib/go/thrift"
    "math"
)

func Do(w http.ResponseWriter, r *http.Request){
    fmt.Println("\n---------------------\n")
    fmt.Println("\ndealing the get methond\n")

    r.ParseForm()
    if r.Method == "GET"{
        //生成请求内容调用thrift客户端的接口
    }else{
        fmt.Println("method is not GET\n")
        fmt.Println("\n---------------------\n")
        return
    }
    fmt.Println("\n---------------------\n")
    return
}
