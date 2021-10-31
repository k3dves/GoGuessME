package main

import (
	"fmt"
	"net"
)

type player struct {
	conn *net.Conn
	nick string
}

func (p *player) sendMsgToPlayer(msg string) {
	barr := []byte(msg)
	_, err := (*p.conn).Write(barr)

	if err != nil {
		fmt.Printf("Error sendding message to player %s, ERROR: %s \n", p.nick, err)
	}

}
