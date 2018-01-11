package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func main() {
	// c, err := redis.Dial("tcp", "192.168.2.163:6379", redis.DialPassword("ciwongrds"))
	rpool := redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", "192.168.2.163:6379", redis.DialPassword("ciwongrds"))
	}, 10)
	defer rpool.Close()

	c := rpool.Get()

	d, e := redis.StringMap(c.Do("HGETALL", "cw:gateway:token:01c48d98fa00440e92930960ca44ddba959f0ffd"))
	fmt.Println(d, e)

	// _, err := c.Do("SET", "gokey", "govalue", "EX", "5")
	_, err := c.Do("SET", "gokey", "govalue")
	if err != nil {
		fmt.Println(err)
	}

	v, err := redis.String(c.Do("GET", "gokey"))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(v)
	}
}
