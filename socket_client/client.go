package main

import (
	"log"
	"net"
	"time"
)

func establishConn(i int) net.Conn {
	// conn, err := net.Dial("tcp", ":8888")
	conn, err := net.DialTimeout("tcp", ":8888", 1*time.Nanosecond)
	if err != nil {
		log.Printf("%d: dial error: %s", i, err)
		return nil
	}
	if es := conn.SetWriteDeadline(time.Now().Add(1 * time.Nanosecond)); es != nil {
		log.Println(es)
		return nil
	}
	if _, e := conn.Write([]byte("connect to server ok")); e != nil {
		log.Println(e)
		return nil
	}
	log.Println(i, ":connect to server ok")
	return conn
}

func main() {
	var sl []net.Conn
	for i := 1; i < 10; i++ {
		conn := establishConn(i)
		if conn != nil {
			sl = append(sl, conn)
		}
	}
	time.Sleep(time.Second * 10000)
}
