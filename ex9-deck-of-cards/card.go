//go:generate stringer -type=Suit,Rank

package deck

import "fmt"

// Suit represents the suit of the card
type Suit uint8

// Spade is 0 and increments down
const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker // special case
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

// Rank represents the value of the card
type Rank uint8

// Ace is 1 and increments down
const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minRank = Ace
	maxRank = King
)

// Card has a Suit and a Rank to represent a card
type Card struct {
	Suit
	Rank
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

// New returns a deck of Cards
func New() []Card {
	var cards []Card
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}
	return cards
}

func main() {
	fmt.Println("vim-go")
}
