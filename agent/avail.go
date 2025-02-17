package agent

import (
	"fmt"
	"math"

	"solitaire/deck"
	"solitaire/game"
)

func (agent *Agent) findAvailCards() []*deck.Card {
	faceUp := agent.game.Avail
	faceDn := agent.game.Deck

	cap := int(math.Ceil(float64(len(faceUp) + len(faceDn)) / float64(game.NFLIP)))
	cards := make([]*deck.Card, 0, cap)

	if len(faceUp) > 0 {
		// Append current faceup card
		cards = append(cards, &faceUp[len(faceUp) - 1])
	}

	if len(faceDn) > 0 {
		// Get every third card from face-down pile
		for i := 2; i < len(faceDn); i += game.NFLIP {
			cards = append(cards, &faceDn[i])
		}

		if len(faceDn) % game.NFLIP != 0 {
			cards = append(cards, &faceDn[len(faceDn) - 1])
		}
	}

	if len(faceUp) > 0 {
		// Get every third card from face-up pile
		for i := 2; i < len(faceUp) - 1; i += game.NFLIP {
			// len(faceUp) - 1 so we skip the originally added card
			cards = append(cards, &faceUp[i])
		}

		if len(faceUp) % game.NFLIP != 0 {
			// If nFaceUp not divisible by 3
			if len(faceDn) > 0 { 
				// and there are face-down cards, go back through face-down pile
				for i := (game.NFLIP - 1) - (len(faceUp) % game.NFLIP); i < len(faceDn); i += game.NFLIP {
					cards = append(cards, &faceDn[i])
				}

				if (len(faceDn) + len(faceUp)) % 3 != 0 && len(faceDn) % 3 != 0 {
					cards = append(cards, &faceDn[len(faceDn) - 1])
				}
			}

		}
	}

	return cards
}



func TestFindAvailCards() {
	fmt.Println("Hello world!")

	deck := deck.NewDeck()
	game := game.NewGame(deck)

	for _ = range(4) {
		game.Flip()
		fmt.Println("Flipped!")
	}
	game.MoveAvailToTop()
	for _ = range(5) {
		game.Flip()
		fmt.Println("Flipped!")
	}

	game.Display(false)

	strat := NullStrategy{}
	agent,_ := NewAgent(game, strat)

	for _,cardPtr := range agent.findAvailCards() {
		fmt.Print(*cardPtr, ", ")
	}
	fmt.Print("\n")
}
