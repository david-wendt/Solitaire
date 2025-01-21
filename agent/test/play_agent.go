package main

import (
	"solitaire/agent"
	"solitaire/deck"
	"solitaire/game"
	"solitaire/ioutils"
)

func main() {

	deck := deck.NewDeck()
	deck.Shuffle()

	game := game.NewGame(deck)

	// Idea: Same initialization, but priority strategy. 
	// Deterministic, and move down in order of priority.
	// Might be hard to code dynamically, easier to hardcode the priority.
	strategy := agent.ProbabilisticStrategy{
		PFlip: 0.001, // Warning! If this is zero, can get stuck in infinite loop
		PTableau: 0.099,
		PAvail: 0.89, // This does not do what I want, since
		PToTop: 0.1, // Avail includes AvailToTop.
		PFromTop: 0.0,
	}

	agent,err := agent.NewAgent(game, strategy)
	if err != nil {
		panic(err)
	}
	
	game.Display(true)
	for {
		input := ioutils.Input("Continue? <Enter> if so, anything else to quit.")
		if input != "" { break }

		agent.Act(true)
		game.Display(true)
	}
}