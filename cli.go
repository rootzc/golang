package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"myhttp1.1/gen-go/demo/rpc"
	"net"
	"os"
	"time"
    "net/http"
    "strings"
    "myhttp1.1/redisopt"
    //"strconv"
)

const (
    Netaddrmaster = "127.0.0.1"
    portmaster = "19090"
    Netaddrslave = "127.0.0.1"
    portslave = "19091"
)
func CreateThriftClient(netaddr,port string)(client *rpc.DoRedisClient){
    
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(net.JoinHostPort(netaddr, port))
    redisopt.Logmapcli["infolog"].Println("thrift server server on ",netaddr,port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
        redisopt.Logmapcli["errlog"].Println("err on resolving server addr")
		return nil
	}

	useTransport := transportFactory.GetTransport(transport)
	client = rpc.NewDoRedisClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "error opening socket to ", netaddr,port, err)
        redisopt.Logmapcli["errlog"].Println("thrift cannot open socket server on ",netaddr,port)
		return nil
	}


	//defer transport.Close()
    return  client
	//调用函数

}

func call(w http.ResponseWriter,r *http.Request) {
    //thrift框架代码

	//    for i:=0;i<1000;i++{
	starttime := currentTimeMillis()
    
    client := CreateThriftClient(Netaddrmaster,portmaster)
    if client == nil{
        redisopt.Logmapcli["errlog"].Println("[alert]master is down!!!!")
        fmt.Println("master server is down")
        client = CreateThriftClient(Netaddrslave,portslave)
        if client == nil{
            fmt.Println("slave server is down")
            redisopt.Logmapcli["errlog"].Println("[alert]slave is down!!!!")
            os.Exit(1)
        }
    }

    
    //解析请求
    r.ParseForm()
    fmt.Println(r.Form)
    redisopt.Logmapcli["infolog"].Println("httpserver recv ",r.Form)
    command,ok := r.Form["command"]
    if !ok{
        fmt.Println("cannot find the command")
        redisopt.Logmapcli["errlog"].Println("httpserver cannot get  the command from url")
        return
    }
    if len(command)!=1{
        fmt.Println("command only be one")
        redisopt.Logmapcli["errlog"].Println("command filed is morethan one")
        return
    }
    //目前使用string传输，为了做到范型应该转化为interface{}
    //argv := make([]string,0)
    formargv,ok:=r.Form["argv"]
    if !ok{
        fmt.Println("no agrv")
        redisopt.Logmapcli["errlog"].Println("argc is wrong")
        return
    }
    argv:=strings.Split(formargv[0],",")
    //val  := "argv"
    //for i:=0;i<len(r.Form);i++{
    //        
    //        if v,ok:=r.Form[val];ok{
    //            argv = append(argv,v)
    //        }else{
    //        
    //        }
    //        i++
    //        val += strconv.Itoa(i)
    //    }
    //
    //调用嗯thrift接口
	//////////////////////////////////////
    //构建请求报文
    argc := int32(len(argv))

	req := &rpc.Request{Command: command[0], Argc: argc, Argv:argv}
    redisopt.Logmapcli["infolog"].Println("httpserver Request:",req)
	fmt.Println(req)
    res, err := client.Do(req)
	//fmt.Printf("res is %v",res)
	switch {
	case err != nil:
		fmt.Println("err != nil \n", err)
		goto out
		//fallthrough
	case res == nil:
		fmt.Println("res == nil")
		goto out
		//fallthrough
	default:
		switch {
		case res.State == true:
			fmt.Printf("\n-----------\nsucess ! res \n:%v\n--------\n", res.Reslist)
            redisopt.Logmapcli["infolog"].Println("OK : httpserver recv",res.Reslist)
		default:
			fmt.Println("connnect ok but failed State is :", res.State)
			goto out
		}

	}

out:
	////////////////////////////////////
    //结果写入页面
    fmt.Fprintf(w,res.Reslist[0])
	endtime := currentTimeMillis()
	fmt.Println("Pragram exit. time->", endtime, starttime, (endtime - starttime))

}
func main() {
	//        paramMap := make(map[string]string)
    //http

    http.HandleFunc("/",call)
    err:= http.ListenAndServe(":9090",nil)
    if err!=nil{
        fmt.Println("http.ListenAndServe ocurr")
        return
    }
    redisopt.Logmapcli["infolog"].Println("httpserver server on :9090")

}

func currentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}
