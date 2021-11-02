package main

import (
	"fmt"
	"strings"
)

func getAllPlayerName(players map[*player]string) string {
	var arr []string
	fmt.Println(len(players))
	for p := range players {
		d(p.name)
		arr = append(arr, p.name)
	}
	fmt.Print(arr)
	fmt.Print(len(arr))
	return "Players : [" + strings.Join(arr, ", ") + "]\n"
}
