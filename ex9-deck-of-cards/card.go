//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

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
func New(opts ...func([]Card) []Card) []Card {
	var cards []Card
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}
	for _, opt := range opts {
		cards = opt(cards)
	}
	return cards
}

// Sort takes in a custom "less" function
func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

// DefaultSort sorts a deck of cards with the default Less option
func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

// Less takes in a slice of cards returns a less function
func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) < absRank(cards[j])
	}
}

// absRank returns an absolute rank of a card
func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

// Shuffle an array of Cards using a random permutation
func Shuffle(cards []Card) []Card {
	ret := make([]Card, len(cards))
	r := rand.New(rand.NewSource(time.Now().Unix()))
	perm := r.Perm(len(cards))
	// perm = [0,1,4,2,3...]
	for i, j := range perm {
		ret[i] = cards[j]
	}
	return ret
}
