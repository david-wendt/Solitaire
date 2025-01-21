package deck 

import (
	"fmt"
	"math/rand"
)

type RankT byte 
type SuitT byte

type Card struct {
	Rank RankT
	Suit SuitT
}

func NewCard(rank int, suit int) Card {
	return Card{Rank: RankT(rank), Suit: SuitT(suit)}
}

func (c Card) String() string {
    return fmt.Sprintf("%v%v", c.Rank, c.Suit)
}

func (c Card) Color() byte {
	return byte(c.Suit) % 2
}

const (
	Spades SuitT = iota
	Hearts
	Clubs
	Diamonds
)

// // Light mode (dark text, spades and clubs filled)
// var Suits = [...]string{
//     Spades:   "\u2660",
//     Hearts:   "\u2661",
//     Clubs:    "\u2663",
//     Diamonds: "\u2662",
// }

// Dark mode (light text, hearts and diamonds filled)
var Suits = [...]string{
    Spades:   "\u2664",
    Hearts:   "\u2665",
    Clubs:    "\u2667",
    Diamonds: "\u2666",
}

const NSuits = len(Suits)

func (s SuitT) String() string {
    if int(s) < len(Suits) {
        return Suits[s]
    }
    return fmt.Sprintf("Invalid Suit: %v", byte(s))
}

const (
	Ace RankT = iota 
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

var Ranks = [...]string{
    Ace:   "A",
    Two:   "2",
    Three: "3",
    Four:  "4",
    Five:  "5",
    Six:   "6",
    Seven: "7",
    Eight: "8",
    Nine:  "9",
    Ten:   "10",// "\u2491" for 10. char // or "X"? or "T"?
    Jack:  "J",
    Queen: "Q",
    King:  "K",
}

const SuitSize = len(Ranks)

func (r RankT) String() string {
    if int(r) < SuitSize {
        return Ranks[r]
    }
    return fmt.Sprintf("Invalid Rank: %v", byte(r))
}

type Deck [52]Card 

func NewDeck() (d Deck) {
	for iSuit := 0; iSuit < NSuits; iSuit++ {
		for iRank := 0; iRank < SuitSize; iRank++ {
			d[iSuit*SuitSize + iRank] = NewCard(iRank, iSuit)
		}
	}
	return
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

func CanPlace(card, dest Card) bool {
	return dest.Rank == card.Rank + 1 && card.Color() != dest.Color()
}