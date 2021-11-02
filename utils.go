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
