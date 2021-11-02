package main

import (
	"fmt"
	"net"
)

type server struct {
	players          map[string]*player
	listener         net.Listener
	broadcastChannel chan message
	register         chan *player
	deregister       chan *player
	votes            map[string]int
	isVotingEnabled  bool
}

type message struct {
	player *player
	nick   string
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
			s.players[player.nick] = player
			s.broadcastChannel <- message{player: player, nick: player.name, event: "JOIN"}
			fmt.Printf("<%s> joined the server\n", player.name)

		case player := <-s.deregister:
			delete(s.players, player.nick)
			delete(s.votes, player.nick)
			s.broadcastChannel <- message{player: player, nick: player.name, event: "LEFT"}
			(*player.conn).Close()

		case msg := <-s.broadcastChannel:
			s.broadcast(&msg)

		}
	}
}

func (s *server) broadcast(msg *message) {
	for _, player := range s.players {
		if msg.event != "TEXT" || msg.player.nick != player.nick {
			player.sendMsgToPlayer(msg)
		}
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
	s.players = make(map[string]*player)
	s.votes = make(map[string]int)

	if err != nil {
		fmt.Printf("Error starting the server %s", err)
	}

	fmt.Printf("Starting Game Server\n")
	s.run()
}
