package main

import (
	"fmt"

	"github.com/acharyab/gophercises/ex11-blackjack-ai/blackjack"
)

func main() {
	game := blackjack.New()
	winnings := game.Play(blackjack.HumanAI())
	fmt.Println(winnings)
}
