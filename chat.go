package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main()  {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil{
		log.Fatal(err)
	}

	go brodcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		go handleConn(conn)
	}
}

type client chan string

var(
	entering = make(chan client)
	leaving = make(chan client)
	messages = make(chan string)
)

func brodcaster()  {
	clients := make(map[client]bool)

	for{
		select {
		case mes := <-messages:
			for cli := range clients{
				cli <- mes
			}
		case new_cli := <-entering:
			clients[new_cli] = true
		case leaving_cli := <-leaving:
			delete(clients, leaving_cli)
			close(leaving_cli)
		}
	}
}

func handleConn(conn net.Conn){
	ch := make(chan string)
	go clientWriter(conn, ch)

	identity := conn.RemoteAddr().String()
	ch <- "You are " + identity
	messages <- identity + " arrived."

	entering <- ch
	input := bufio.NewScanner(conn)
	for input.Scan(){
		messages<- identity + ": "+ input.Text()
	}

	leaving <- ch
	messages <- identity + "left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch chan string){
	for message := range ch{
		fmt.Fprintln(conn, message)
	}
}