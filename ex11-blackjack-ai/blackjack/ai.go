package blackjack

import (
	"fmt"

	deck "github.com/acharyab/gophercises/ex9-deck-of-cards"
)

// AI has methods required to implement a dealer or blackjack AI
type AI interface {
	Bet(shuffled bool) int
	Results(hands [][]deck.Card, dealer []deck.Card)
	Play(hand []deck.Card, dealer deck.Card) Move
}

type humanAI struct{}

// HumanAI returns a humanAI struct
func HumanAI() AI {
	return humanAI{}
}

func (ai humanAI) Bet(shuffled bool) int {
	if shuffled {
		fmt.Println("The deck was just shuffled")
	}
	fmt.Println("What would you like to bet?")
	var bet int
	fmt.Scanf("%d\n", &bet)
	return bet
}

type dealerAI struct{}

func (ai dealerAI) Bet(shuffled bool) int {
	// noop
	return 1
}

func (ai dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	// If dealer score <= 16, we hit
	// If dealer has a soft 17, then we hit.
	dScore := Score(hand...)
	if dScore <= 16 || dScore == 17 && Soft(hand...) {
		return MoveHit
	}
	return MoveStand
}

func (ai dealerAI) Results(hands [][]deck.Card, dealer []deck.Card) {
	// noop
}

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		fmt.Println("Player:", hand)
		fmt.Println("Dealer:", dealer)
		fmt.Println("What will you do? (h)it, (s)tand, (d)ouble, s(p)lit")
		var input string
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		case "d":
			return MoveDouble
		case "p":
			return MoveSplit
		default:
			fmt.Println("Invalid option: ", input)
		}
	}
}

func (ai humanAI) Results(hands [][]deck.Card, dealer []deck.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:")
	for _, h := range hands {
		fmt.Println(" ", h)
	}
	fmt.Println("Dealer:", dealer)
}

func main() {
	fmt.Println("vim-go")
}
