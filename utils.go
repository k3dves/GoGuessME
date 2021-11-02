package main

import (
	"fmt"
	"strings"
)

func getAllPlayerName(players map[string]*player) string {
	var arr []string
	fmt.Println(len(players))
	for nick := range players {
		d(nick)
		arr = append(arr, nick)
	}
	fmt.Print(arr)
	fmt.Print(len(arr))
	return "Players : [" + strings.Join(arr, ", ") + "]\n"
}

func getVotesAsString(votes map[string]int) string {
	res := "["
	for name, count := range votes {
		res += fmt.Sprintf("%s:%d ,", name, count)
	}
	res += "]\n"
	return res

}

func votePlayer(server *server, voter string, voted string) string {
	if !server.players[voter].canVote {
		return "You already voted!! \n"
	}
	if !server.players[voter].alive {
		return "You're dead can't vote!!\n"
	}
	votes := server.votes
	if val, ok := votes[voted]; ok {
		votes[voted] = val + 1
		server.players[voter].canVote = false
		return "Voted Player : " + voted + "\n"
	}

	return "Error , Player not present!!\n"
}

func initVotes(votes *map[string]int) {
	for key := range *votes {
		(*votes)[key] = 0
	}
}

func enableVotingForEachPlayer(players map[string]*player) {
	for _, player := range players {
		player.canVote = true
	}
}
