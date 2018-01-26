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
			d := make([]byte, 512)
			var n int
			for {
				time.Sleep(time.Second * 1)
				if _, es := c.Write([]byte("hello!")); es != nil {
					log.Println(es)
				}
				var er error
				if n, er = c.Read(d); er != nil {
					log.Println(er)
				}
				log.Println(string(d[:n]))
			}
		}()
		i++
		log.Printf("%d: accept a new connection\n", i)
		log.Println("local: ", c.LocalAddr().String(), ", remote: ", c.RemoteAddr().String())
	}
}
