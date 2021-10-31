package main

import (
	"fmt"
	"net"
)

type player struct {
	conn *net.Conn
	nick string
	name string
}

func (p *player) sendMsgToPlayer(msg message) {
	var err error
	switch msg.event {
	case "JOIN":
		barr := []byte("> " + msg.nick + " joined the server\n")
		_, err = (*p.conn).Write(barr)

	case "LEFT":
		barr := []byte("< " + msg.nick + " left the server\n")
		_, err = (*p.conn).Write(barr)

	case "TEXT":
		barr := []byte("<" + msg.nick + "> " + msg.msg)
		_, err = (*p.conn).Write(barr)

	case "SERVER":
		barr := []byte("<server> " + msg.msg)
		_, err = (*p.conn).Write(barr)
	}

	if err != nil {
		fmt.Printf("Error sendding message to player %s, ERROR: %s \n", p.nick, err)
	}

}
