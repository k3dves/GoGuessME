package main

import (
	"fmt"
	"net"
)

type player struct {
	conn *net.Conn
	nick string
}

func (p *player) sendMsgToPlayer(msg message) {
	var err error
	switch msg.event {
	case "JOIN":
		barr := []byte("> " + msg.name + " joined the server\n")
		_, err = (*p.conn).Write(barr)

	case "LEFT":
		barr := []byte("< " + msg.name + " left the server\n")
		_, err = (*p.conn).Write(barr)

	case "TEXT":
		barr := []byte("<" + msg.name + "> " + msg.msg)
		_, err = (*p.conn).Write(barr)

	}

	if err != nil {
		fmt.Printf("Error sendding message to player %s, ERROR: %s \n", p.nick, err)
	}

}
