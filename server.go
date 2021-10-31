package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

type server struct {
	players          map[*player]string
	listener         net.Listener
	broadcastChannel chan message
	register         chan *player
	deregister       chan *player
}

type message struct {
	player *player
	name   string
	msg    string
	event  string
}

func (s *server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting client connection %s\n", err)
		}
		go handleNewPlayer(s, &conn)
	}
}

func (s *server) run() {
	go s.listen()

	for {

		select {
		case player := <-s.register:
			s.players[player] = player.nick
			s.broadcastChannel <- message{player: player, name: player.nick, event: "JOIN"}
			fmt.Printf("<%s> joined the server\n", player.nick)

		case player := <-s.deregister:
			delete(s.players, player)
			s.broadcastChannel <- message{player: player, name: player.nick, event: "LEFT"}
			(*player.conn).Close()

		case msg := <-s.broadcastChannel:
			s.broadcast(msg)

		}
	}
}

func (s *server) broadcast(msg message) {
	for player := range s.players {
		if msg.event != "TEXT" || msg.player != player {
			player.sendMsgToPlayer(msg)
		}
	}

}
func handleNewPlayer(s *server, c *net.Conn) {
	p := &player{}
	p.conn = c
	(*p.conn).Write([]byte(PLAYER_GREET))
	reader := bufio.NewReader(*p.conn)
	nick, err := reader.ReadString('\n')

	if err != nil {
		fmt.Printf("Error getting player nick : %s\n", err)
	}
	p.nick = strings.TrimSuffix(nick, "\n")

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
		s.broadcastChannel <- message{player: p, name: p.nick, msg: msg, event: "TEXT"}
	}

}

//StartGameServer -> entrypoint for server.go
func StartGameServer() {
	s := &server{}
	var err error
	s.listener, err = net.Listen("tcp", ":6060")
	s.register = make(chan *player)
	s.deregister = make(chan *player)
	s.broadcastChannel = make(chan message, BRODCAST_CHAN_SIZE)
	s.players = make(map[*player]string)

	if err != nil {
		fmt.Printf("Error starting the server %s", err)
	}

	fmt.Printf("Starting Game Server\n")
	s.run()
}
