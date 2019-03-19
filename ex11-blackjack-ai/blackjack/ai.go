package blackjack

import (
	"fmt"

	deck "github.com/acharyab/gophercises/ex9-deck-of-cards"
)

type AI interface {
	Bet() int
	Results(hand [][]deck.Card, dealer []deck.Card)
	Play(hand []deck.Card, dealer deck.Card) Move
}

type humanAI struct{}

func HumanAI() AI {
	return humanAI{}
}

func (ai humanAI) Bet() int {
	return 1
}

type dealerAI struct{}

func (ai dealerAI) Bet() int {
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

func (ai dealerAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	// noop
}

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		fmt.Println("Player:", hand)
		fmt.Println("Dealer:", dealer)
		fmt.Println("What will you do? (h)it, (s)tand")
		var input string
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		default:
			fmt.Println("Invalid option: ", input)
		}
	}
}

func (ai humanAI) Results(hand [][]deck.Card, dealer []deck.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:", hand)
	fmt.Println("Dealer:", dealer)
}

func main() {
	fmt.Println("vim-go")
}
