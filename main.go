package main

import (
	"io"
	"log"
	"net"
)

func handleClient(conn net.Conn) {
	defer func() {
		log.Printf("%v logout", conn.RemoteAddr())
		conn.Close()
	}()

	log.Printf("Connection from %v", conn.RemoteAddr())

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if n > 0 {
			log.Printf("%s", buf[:n])
		}
		if err != nil {
			if err != io.EOF {
				log.Print(err)
			}
			return
		}

	}
}

func main() {
	ln, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}

		go handleClient(conn)
	}
}
