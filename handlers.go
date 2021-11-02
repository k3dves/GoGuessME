package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

func handleNewPlayer(s *server, c *net.Conn) {
	player := &player{}
	player.conn = c

	(*player.conn).Write([]byte(PLAYER_GREET))
	reader := bufio.NewReader(*player.conn)
	msg, err := reader.ReadString('\n')

	if err != nil {
		fmt.Printf("Error getting player nick : %s\n", err)
		return
	}
	msg = strings.TrimSuffix(msg, "\n")
	arr := strings.Split(msg, ":")
	player.name = arr[0]
	//Todo only allow unique nicks
	player.nick = strings.ToUpper(arr[1])
	s.votes[player.nick] = 0
	s.register <- player
	for {
		//Blocks till it gets a string
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			//player left
			s.deregister <- player
			break
		}

		if err != nil {
			fmt.Printf("Error reading player: %s message. Error: %s \n", player.nick, err)
			break
		}
		//Keep receiving the mgs and pass it to the handler
		messageHandler(s, player, text)
	}

}

func messageHandler(server *server, player *player, text string) {
	//cerate a new message object
	message := message{player: player, nick: player.nick}

	if strings.HasPrefix(text, CMD_IDENTIFIER) {
		text = strings.TrimSuffix(text, "\n")
		text = strings.TrimSuffix(text, " ")
		cmd := strings.TrimPrefix(text, CMD_IDENTIFIER)
		cmd = strings.ToUpper(cmd)
		d("command reveived " + cmd)
		message.event = "CMD"
		message.msg = cmd
		commandHandler(server, &message)
	} else {
		message.event = "TEXT"
		message.msg = text
		server.broadcastChannel <- message
	}
}
func commandHandler(server *server, message *message) {
	commandMap := commandParser(message)
	switch commandMap["cmd"] {
	case "SHOW":
		playerNames := getAllPlayerName(server.players)
		message.msg = playerNames
	case "VOTE":
		voteResult := votePlayer(server.votes, commandMap["option"])
		fmt.Print(server.votes)
		message.msg = voteResult
	default:
		message.msg = "Invalid command " + message.msg + "\n"

	}
	message.event = "SERVER"
	message.player.sendMsgToPlayer(message)

}

func commandParser(message *message) map[string]string {
	parsed := make(map[string]string)
	arr := strings.Split(message.msg, " ")
	parsed["cmd"] = arr[0]
	if len(arr) > 1 {
		parsed["option"] = arr[1]
	}

	return parsed
}

func d(m string) {
	fmt.Println(m)
}
