package redisopt

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

//全局的redis对象
//redis命令表
//var RedisDB  RedisObj


var CommandTable  = make(map[string]Command)

const (
	net = ":6379"
)

//包的初始化函数
func init() {
	//RedisDB.Init(":6379",3)

	var get Get //:= Get{nil}
	var set Set //:= Set{nil}

	CommandTable["GET"] = get

	CommandTable["SET"] = set

}

//////////////////////////////////////////////////////////
//redis数据库的抽象

type RedisDB struct {
	netaddr string
	pool    *redis.Pool
}

func (this *RedisDB) Close() {
	this.pool.Close()
}

func NewRedisDB(redisaddr string, idle int) (db *RedisDB) {
	db = &RedisDB{netaddr: redisaddr, pool: redis.NewPool(func() (redis.Conn, error) { return redis.Dial("tcp", redisaddr) }, 3)}
	return db
}

//初始化redis网络地址以及链接池
//从链接池得到一个链接
func (this *RedisDB) GetConn() (conn redis.Conn) {
	conn = this.pool.Get()
	return conn
}

var (
	errragc = errors.New("argc not right")
	errtype = errors.New("argv type not right")
)

///////////////////////////////////////////////////////////////////////////////

//redis命令
type Command interface{}

//
//子类Get方法
type Get struct {
	db *RedisDB
}

func (this *Get) Run(argc int, argv string) (res string, err error) {
	if argc != 1 {
		return "", errragc
	}
	//

	this.db = NewRedisDB(net, 3)
	conn := this.db.GetConn()
	defer conn.Close()

	//
	res, err = redis.String(conn.Do("GET", argv))
	if err != nil {
		//此处添加日志
		fmt.Println("Get.Run err")

	}
	//fmt.Println(res)
	return
}

//子类set方法
type Set struct {
	db *RedisDB
}

func (this *Set) Run(argc int, argv []interface{}) (res string, err error) {
	if argc != 2 {
		return "false", errragc
	}
	key, ok := argv[0].(string)
	if !ok {
		return "false", errtype
	}
	//val,ok:= argv[1].(string)
	//if !ok {
	//    err = errors.New("cannot get the val")
	//    return
	//}

	//
	this.db = NewRedisDB(net, 3)
	conn := this.db.GetConn()
	defer conn.Close()
	//
	res, err = redis.String(conn.Do("SET", key, argv[1]))
	if err != nil {
		//记录日志
		fmt.Println("Set.Run() err:", err)

	}
	//记录日志
	fmt.Println(res)
	return res,err
}

