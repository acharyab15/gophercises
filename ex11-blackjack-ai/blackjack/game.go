package blackjack

import (
	"errors"
	"fmt"

	deck "github.com/acharyab/gophercises/ex9-deck-of-cards"
)

// State represents the state of a Turn
type state int8

// statePlayerTurn starts at 0 and increments down
const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

// New returns a new Game
func New(opts Options) Game {
	g := Game{
		state:    statePlayerTurn,
		dealerAI: dealerAI{},
		balance:  0,
	}
	if opts.Decks == 0 {
		opts.Decks = 3
	}
	if opts.Hands == 0 {
		opts.Hands = 100
	}
	if opts.BlackjackPayout == 0.0 {
		opts.BlackjackPayout = 1.5
	}
	g.nDecks = opts.Decks
	g.nHands = opts.Hands
	g.blackjackPayout = opts.BlackjackPayout
	return g
}

// Options for the game that a user can specify
type Options struct {
	Decks           int
	Hands           int
	BlackjackPayout float64
}

// Game defines the individual fields required for a blackjack game
type Game struct {
	// unexported fields
	nDecks          int
	nHands          int
	blackjackPayout float64

	deck  []deck.Card
	state state

	player    []deck.Card
	playerBet int
	balance   int

	dealer   []deck.Card
	dealerAI AI
}

// CurrentPlayer returns the current player's hand
func (g *Game) currentHand() *[]deck.Card {
	switch g.state {
	case statePlayerTurn:
		return &g.player
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("it isn't currently any player's turn")
	}

}

// deal two cards to the players
func deal(g *Game) {
	g.player = make([]deck.Card, 0, 5)
	g.dealer = make([]deck.Card, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, g.deck = draw(g.deck)
		g.player = append(g.player, card)
		card, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, card)
	}
	g.state = statePlayerTurn
}

// bet allows the player to place bets
func bet(g *Game, ai AI, shuffled bool) {
	bet := ai.Bet(shuffled)
	g.playerBet = bet
}

// Play the game a specified number of times
// First, it's the player turn and then the dealer turn
func (g *Game) Play(ai AI) int {
	g.deck = nil
	min := 52 * g.nDecks / 3 // 3 is arbitrary
	for i := 0; i < g.nHands; i++ {
		shuffled := false
		if len(g.deck) < min {
			g.deck = deck.New(deck.Deck(g.nDecks), deck.Shuffle)
			shuffled = true
		}
		bet(g, ai, shuffled)
		deal(g)
		if Blackjack(g.dealer...) {
			endHand(g, ai)
			continue
		}
		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(g.player))
			copy(hand, g.player)
			move := ai.Play(hand, g.dealer[0])
			err := move(g)
			switch err {
			case errBust:
				MoveStand(g)
			case nil:
				//noop
			default:
				panic(err)
			}
		}
		for g.state == stateDealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, g.dealer[0])
			move(g)
		}
		endHand(g, ai)
	}
	return g.balance
}

var (
	errBust = errors.New("hand score exceeded 21")
)

// Move is either a MoveHit or a MoveStand
type Move func(*Game) error

// MoveHit draws a card from the deck and
// and appends it to the hand
func MoveHit(g *Game) error {
	hand := g.currentHand()
	var card deck.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		MoveStand(g)
	}
	return nil
}

// MoveDouble allows the user to double down
func MoveDouble(g *Game) error {
	if len(g.player) != 2 {
		return errors.New("can only double on a hand with 2 cards")
	}
	g.playerBet *= 2
	MoveHit(g)
	return MoveStand(g)
}

// MoveStand changes game state to the next state
// either from player -> dealer or dealer -> handOver
func MoveStand(g *Game) error {
	g.state++
	return nil
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

// endHand performs some end of game logic
func endHand(g *Game, ai AI) {
	pScore, dScore := Score(g.player...), Score(g.dealer...)
	pBlackjack, dBlackjack := Blackjack(g.player...), Blackjack(g.dealer...)
	winnings := g.playerBet
	// TODO(acharyab): Figure out winnings and add/subtract them
	switch {
	case pBlackjack && dBlackjack:
		winnings = 0
	case dBlackjack:
		winnings = -winnings
	case pBlackjack:
		winnings = int(float64(winnings) * g.blackjackPayout)
	case pScore > 21:
		winnings = -winnings
		g.balance--
	case dScore > 21:
		// win
	case pScore > dScore:
		// win
	case dScore > pScore:
		winnings = -winnings
	case dScore == pScore:
		winnings = 0
	}
	g.balance += winnings
	fmt.Println()
	ai.Results([][]deck.Card{g.player}, g.dealer)
	g.player = nil
	g.dealer = nil
}

// Score takes in a hand of cards and returns the best blackjack score possible
func Score(hand ...deck.Card) int {
	minScore := minScore(hand...)
	if minScore > 11 {
		return minScore
	}
	for _, c := range hand {
		if c.Rank == deck.Ace {
			// ace is currently worth 1, making it 11
			return minScore + 10
		}
	}
	return minScore

}

// minScore represents the minimum score of the hand
func minScore(hand ...deck.Card) int {
	score := 0
	for _, card := range hand {
		score += min(10, int(card.Rank))
	}
	return score
}

// Blackjack returns true if a hand is a blackjack
func Blackjack(hand ...deck.Card) bool {
	if len(hand) == 2 && Score(hand...) == 21 {
		return true
	}
	return false
}

// Soft returns if the score of a hand is a soft score
// i.e. if an ace is being counted as 11 points
func Soft(hand ...deck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)
	return minScore != score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
