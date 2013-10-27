package main
 
import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)
 
type Message struct {
	conn *net.TCPConn
	msg  string
}
 
var (
	messages     = make(chan Message)
	connections  = make(map[string]*net.TCPConn)
	connectionsM sync.RWMutex
)
 
func handleClient(conn *net.TCPConn) {
	defer func() {
		log.Printf("Disconnect from %v", conn.RemoteAddr())
		conn.Close()
		connectionsM.Lock()
		delete(connections, conn.RemoteAddr().String())
		connectionsM.Unlock()
	}()
 
	log.Printf("Connect from %v", conn.RemoteAddr())
 
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			messages <- Message{conn, string(buf[:n])}
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("Error from %v: %s", conn.RemoteAddr(), err)
			}
			return
		}
	}
}
 
func main() {
	go func() {
		for m := range messages {
			connectionsM.RLock()
			for _, conn := range connections {
				if m.conn != conn {
					s := fmt.Sprintf("%s: %s", m.conn.RemoteAddr(), m.msg)
					_, err := conn.Write([]byte(s))
					if err != nil {
						log.Print(err)
					}
				}
			}
			connectionsM.RUnlock()
		}
	}()
 
	ln, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatal(err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}
 
		conn := c.(*net.TCPConn)
		connectionsM.Lock()
		connections[conn.RemoteAddr().String()] = conn
		connectionsM.Unlock()
		go handleClient(conn)
	}
}