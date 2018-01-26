package main

import (
	"log"
	"net"
	"strconv"
	"time"
)

func establishConn(i int) net.Conn {
	conn, err := net.Dial("tcp", ":8888")
	// conn, err := net.DialTimeout("tcp", ":8888", 1*time.Nanosecond)
	if err != nil {
		log.Printf("%d: dial error: %s", i, err)
		return nil
	}
	// if es := conn.SetWriteDeadline(time.Now().Add(1 * time.Nanosecond)); es != nil {
	// 	log.Println(es)
	// 	return nil
	// }
	// if _, e := conn.Write([]byte("connect to server ok")); e != nil {
	// 	log.Println(e)
	// 	return nil
	// }

	go func() {
		d := make([]byte, 512) // 一次读取512byte数据
		var n int
		for {
			time.Sleep(time.Second * 1)
			if _, es := conn.Write([]byte("消息 msg from client connection: " + strconv.Itoa(i))); es != nil {
				log.Println(es)
			}
			var er error
			if n, er = conn.Read(d); er != nil {
				log.Println(er)
			}
			log.Println(strconv.Itoa(i) + ": server msg: " + string(d[:n]))
		}
	}()
	log.Println(i, ":connect to server ok")
	log.Println("local: ", conn.LocalAddr().String(), ", remote: ", conn.RemoteAddr().String())
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
