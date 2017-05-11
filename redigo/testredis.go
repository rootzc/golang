package main

import (
	 //"errors"
	"fmt"
    //"time"
	"github.com/garyburd/redigo/redis"
	_ "strings"
)

func main() {
	/*
	   c,err:=redis.Dial("tcp",":6379")
	   if(err!=nil){
	       fmt.Println("redis connot connect")
	       return
	   }
	   defer c.Close()
	   fmt.Println("sucess connnect redis... ")

	   c.Do("SET","foo",1)
	   exists,_ := redis.Bool(c.Do("exists","foo"))
	   fmt.Println("%v\n",exists)
	*/
	//use pool
	var (
		pool *redis.Pool
		//redisServer := flag.String("redisServer",":6379","")
	)
	//flag.Parse()
	const addr = ":6379"

	pool = redis.NewPool(
		func() (redis.Conn,error){ return redis.Dial("tcp", addr) },
		3)

	defer pool.Close()

    c := pool.Get()
	defer c.Close()
//	res, err := c.Do("GET", "foo")
//	if err != nil {
//		fmt.Println("GET foo wrong")
//		return
//	}
    
	fmt.Println(redis.Int(c.Do("GET","foo")))
}
