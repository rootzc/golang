
namespace go demo.rpc

struct Request{
  1:required string command;
  2:required i32 argc;
  3:optional list<string> argv;
}

enum Errcode{
    no_wrong,
    argc_wrong,
    argv_wrong,
    command_wrong,
    cannot_run
}

struct Result{
  1:required bool state;
  2:optional Errcode err;
  3:optional list<string>reslist;
}

#exception ExceptionReq{
#  1:required i32 code;
#  2:optional string reason;
#}

service DoRedis{
    Result Do(1:Request req);
}
