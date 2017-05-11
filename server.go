package main
import (
    ."myhttp1.1/gen-go/demo/rpc"
    "fmt"
    "git.apache.org/thrift.git/lib/go/thrift"
    //"os"
    "errors"
    "myhttp1.1/redisopt"
    "flag"
    //"strconv"
    //"reflect"
)
const (
    Networkaddr1="127.0.0.1:19090"
    Networkaddr2="127.0.0.1:19091"
)

func CheckAndRun(req *Request,command redisopt.Command)(r *Result,ok bool){

    r = &Result{State:false,Err:Errcode_no_wrong,Reslist:make([]string,0)}
        switch command.(type){
        case redisopt.Get:
            cmd := "GET"
            if req.Argc!=1{
                //err = errors.New("argc is wrong!")
                r.Err = Errcode_argc_wrong
                redisopt.Logmapser["errlog"].Println("request's argc wrong in :",cmd)
            }

            //查找方法表得到对应的方法
            method,ok1 := redisopt.CommandTable[cmd]
            if !ok1{
                fmt.Println("cannot find  from Cbtl:",cmd)
                redisopt.Logmapser["errlog"].Println("cannot find the command :",cmd)
                r.Err = Errcode_command_wrong
            }
            //得到真正的方法
            exe ,ok2:= method.(redisopt.Get)
            if !ok2{
                fmt.Println("cannot imp the  method:",cmd)
                redisopt.Logmapser["errlog"].Println("cannot impl real :GET")        
                r.Err = Errcode_command_wrong 
                
            }
            
            res ,err := exe.Run(int(req.Argc),req.Argv[0])
             
            if err!=nil{
                fmt.Println("err := Run() err!=nil:",cmd)
                redisopt.Logmapser["errlog"].Println("cannot run")
                r.Err = Errcode_cannot_run
                //goto wrong
            }
            r.Reslist = append(r.Reslist,res)
            r.State = true
            ok = true
            fmt.Println("sucess ,get :%v",r.Reslist)
            redisopt.Logmapser["infolog"].Println("sucess ,get :",r.Reslist)
        case redisopt.Set:
            cmd  := "SET"
            if req.Argc!=2{
                //err = errors.New("argc is wrong!")
                r.Err = Errcode_argc_wrong
                redisopt.Logmapser["errlog"].Println("request's argc wrong in :",cmd)
            }

            //查找方法表得到对应的方法
            method,ok1 := redisopt.CommandTable[cmd]
            if !ok1{
                fmt.Println("cannot find GET from Cbtl")
                redisopt.Logmapser["errlog"].Println("cannot find the command :",cmd)
                r.Err = Errcode_command_wrong
            }
            //得到真正的方法
            exe ,ok2:= method.(redisopt.Set)
            if !ok2{
                fmt.Println("cannot find the method",cmd)
                redisopt.Logmapser["errlog"].Println("cannot impl real ",cmd)        
                r.Err = Errcode_command_wrong 
                
            }
            //转化参数为范型
            argv := make([]interface{},0)
            for _,v:= range req.Argv{
                argv = append(argv,v)
            }
            res ,err := exe.Run(int(req.Argc),argv)
             
            if err!=nil{
                fmt.Println("err := Run() err!=nil:",cmd)
                redisopt.Logmapser["errlog"].Println("cannot run")
                r.Err = Errcode_cannot_run
                //goto wrong
            }
            r.Reslist = append(r.Reslist,res)
            r.State = true
            ok = true
            fmt.Println("sucess ,get :%v",r.Reslist)
            redisopt.Logmapser["infolog"].Println("sucess , get Result:",r.Reslist)
            
        } 
        return r,ok
}

var errfailRun = errors.New("command fail")
//thrift接口实现
type Rpcserverimpl struct{

}
func (this *Rpcserverimpl)Do(req *Request)(res *Result,err error){


    //拿到请求的方法，参数
    fmt.Printf("%v\n",req)
    switch {
    case req.Command == "GET":
        var get redisopt.Get
        r,ok := CheckAndRun(req,get)

        if !ok{
            redisopt.Logmapser["errlog"].Println("Get.run(),Fail")
            return r,errfailRun
        }
        res = r
    case  req.Command == "SET":
        var set  redisopt.Set
        r,ok := CheckAndRun(req,set)
        if !ok{
            redisopt.Logmapser["errlog"].Println("Get.run(),Fail")
            return r,errfailRun
        }
        res = r
    }

    return res,nil
}

func CreateThriftServer(netaddr string)(server *thrift.TSimpleServer){

    transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
    protocolfactory := thrift.NewTBinaryProtocolFactoryDefault()
    serverTransport,err := thrift.NewTServerSocket(netaddr)
    if err!=nil{
        fmt.Println("Error!",err)
        redisopt.Logmapser["errlog"].Println("canot create the transport")
        return nil
    }
    handler := &Rpcserverimpl{}
    processor := NewDoRedisProcessor(handler)

    server = thrift.NewTSimpleServer4(processor,serverTransport,transportFactory,protocolfactory)
    fmt.Println("thrift server in ",netaddr)
    redisopt.Logmapser["errlog"].Println("thrift server in ",netaddr)
    return server
}
func main(){
    netaddr := flag.String("addr","127.0.0.1:19090","ip:port")
    flag.Parse()

    server := CreateThriftServer(*netaddr)
    if server == nil{
    //
        redisopt.Logmapser["errlog"].Println("cannot start the thrift server in ",Networkaddr1)
    //    server = CreateThriftServer(Networkaddr1)
    //    if server ==nil{
    //        redisopt.Logmapser["errlog"].Println("cannot start the thrift server in ",Networkaddr2)
    //        os.Exit(1)
        }
    //}
    server.Serve()
}
