package main

import (
	"log"
	"net"
	"time"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("err listen:", err)
		return
	}
	defer l.Close()
	log.Println("listen ok")

	var i int
	for {
		time.Sleep(time.Second * 10)
		c, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			break
		}
		go func() {
			for {
				time.Sleep(time.Second * 1)
				if _, es := c.Write([]byte("accept a new connection")); es != nil {
					log.Println(es)
				}
				var d []byte
				if _, er := c.Read(d); er != nil {
					log.Println(er)
				}
			}
		}()
		i++
		log.Printf("%d: accept a new connection\n", i)
	}
}
