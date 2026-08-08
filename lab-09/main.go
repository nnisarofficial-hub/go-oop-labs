package main

import (
	"fmt"
	"sort"
)

type Player struct {
	Name  string
	Score int
	Rank  int
}

type ByScore []Player

func (b ByScore) Len() int {
	return len(b)
}

func (b ByScore) Less(i, j int) bool {
	return b[i].Score > b[j].Score
}
func (b ByScore) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

type ByName []Player

func (b ByName) Len() int {
	return len(b)
}

func (b ByName) Less(i, j int) bool {
	return b[i].Name < b[j].Name
}
func (b ByName) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func AssignRanks(players []Player) {
	for i := range players {
		players[i].Rank = i + 1
	}

}

func PrintLeaderboard(players []Player) {
	for i := range players {
		fmt.Printf("Rank %d | %-12s| %d pts\n", players[i].Rank, players[i].Name, players[i].Score)
	}
}
func main() {
	players := []Player{
		{Name: "Sara", Score: 9200},
		{Name: "Ali", Score: 9850},
		{Name: "Umar", Score: 9200},
		{Name: "Fatima", Score: 10100},
	}
	sort.Sort(ByScore(players))
	AssignRanks(players)
	PrintLeaderboard(players)
}
