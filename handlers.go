package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

func handleNewPlayer(s *server, c *net.Conn) {
	p := &player{}
	p.conn = c

	(*p.conn).Write([]byte(PLAYER_GREET))
	reader := bufio.NewReader(*p.conn)
	msg, err := reader.ReadString('\n')

	if err != nil {
		fmt.Printf("Error getting player nick : %s\n", err)
		return
	}
	msg = strings.TrimSuffix(msg, "\n")
	arr := strings.Split(msg, ":")
	p.name = arr[0]
	p.nick = arr[1]
	s.register <- p
	for {
		//Blocks till it gets a string
		msg, err := reader.ReadString('\n')
		if err == io.EOF {
			//player left
			s.deregister <- p
			break
		}

		if err != nil {
			fmt.Printf("Error reading player: %s message. Error: %s \n", p.nick, err)
			break
		}
		//Keep receiving the mgs and broadcast to broadcastChannel
		s.broadcastChannel <- message{player: p, nick: p.nick, msg: msg, event: "TEXT"}
	}

}

func d(m string) {
	fmt.Println(m)
}
