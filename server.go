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
	broadcastChannel chan string
	register         chan *player
	deregister       chan *player
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
			s.broadcastChannel <- "> " + player.nick + " joined the server\n"
			fmt.Printf("<%s> joined the server\n", player.nick)

		case player := <-s.deregister:
			delete(s.players, player)
			s.broadcastChannel <- "< " + player.nick + " left the server\n"
			(*player.conn).Close()

		case msg := <-s.broadcastChannel:
			s.broadcast(msg)

		}
	}
}

func (s *server) broadcast(msg string) {
	for player := range s.players {
		player.sendMsgToPlayer(msg)
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
		s.broadcastChannel <- "<" + p.nick + "> " + msg
	}

}

//StartGameServer -> entrypoint for server.go
func StartGameServer() {
	s := &server{}
	var err error
	s.listener, err = net.Listen("tcp", ":6060")
	s.register = make(chan *player)
	s.deregister = make(chan *player)
	s.broadcastChannel = make(chan string, BRODCAST_CHAN_SIZE)
	s.players = make(map[*player]string)

	if err != nil {
		fmt.Printf("Error starting the server %s", err)
	}

	fmt.Printf("Starting Game Server\n")
	s.run()
}
