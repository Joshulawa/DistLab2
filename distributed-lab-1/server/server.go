package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
	_ = fmt.Errorf(err.Error())
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	for {
		conn, _ := ln.Accept()
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	reader := bufio.NewReader(client)
	for {
		msg, _ := reader.ReadString('\n')
		//fmt.Println(msg)
		if len(msg) != 0 {
			msgs <- Message{message: msg, sender: clientid}
		}
		//fmt.println(client, "OK") //Prints to connection the message OK, ie returning a message, think like stdOut
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.
	ln, _ := net.Listen("tcp", *portPtr)
	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	id := 0
	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients channel
			clients[id] = conn //put conn in clients at their id
			// - start to asynchronously handle messages from this client
			go handleClient(conn, id, msgs)
			id++
		case msg := <-msgs:
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
			for id, client := range clients {
				if id != msg.sender {
					fmt.Fprintln(client, "\n", msg.sender, ">  ", msg.message)
				}
			}
		}
	}
}
