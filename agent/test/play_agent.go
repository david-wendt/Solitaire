package main

import (
	"fmt"
	"time"
	"math"

	"solitaire/agent"
	"solitaire/deck"
	"solitaire/game"
	// "solitaire/ioutils"
)

func main() {

	// agent.TestFindAvailCards()
	// return 

	// Idea: Same initialization, but priority strategy. 
	// Deterministic, and move down in order of priority.
	// Might be hard to code dynamically, easier to hardcode the priority.

	// Warning! If PFlip is zero, can get stuck in infinite loop
	// strategy := agent.ProbabilisticStrategy{
	// 	PFlip: 0.000001, 
	// 	PTableau: 1.,
	// 	PAvail: 0.001,
	// 	PToTop: 10000.,
	// 	PFromTop: 0.,
	// }
	strategy := agent.Manual{}
	// verbose := true 
	verbose := false 

	// won := runGame(strategy, true)
	// if won {
	// 	fmt.Println("Congratulations! You won!")
	// } else {
	// 	fmt.Println("You lost :(")
	// }

	nGames := 1
	nTrials := 1

	sumWinRate := 0.0
	sumSqWinRate := 0.0
	for itrial := range nTrials {
		wins := 0
		var sumTime int64 = 0
		var sumSqTime int64 = 0
		for range nGames {
			start := time.Now()
			won := runGame(strategy, verbose)
			if won { wins++ }
			
			if nGames == 1 && nTrials == 1 {
				if won {
					fmt.Println("Congratulations! You won!")
				} else {
					fmt.Println("Oh no! You lost :(")
				}
			}

			t := time.Now()
			elapsed := t.Sub(start)

			nanosecs := elapsed.Nanoseconds()
			sumTime += nanosecs
			sumSqTime += nanosecs*nanosecs
			if sumSqTime > math.MaxInt64 / 2 {panic("Overflow incoming!")}
		}
		// meanTime := sumTime / int64(nGames)
		// meanSqTime := sumSqTime / int64(nGames)
		// stdev := math.Sqrt(float64(meanSqTime - meanTime * meanTime))
		winRate := float64(wins)/float64(nGames)
		fmt.Println("Trial", itrial, "Win rate:", winRate)
		// fmt.Println("Avg time (ns):", meanTime, "Stdev time (ns):", stdev)

		sumWinRate += winRate
		sumSqWinRate += winRate * winRate
	}
	avgWinRate := sumWinRate / float64(nTrials)
	avgSqWinRate := sumSqWinRate / float64(nTrials)
	stdWinRate := math.Sqrt(avgSqWinRate - avgWinRate * avgWinRate)
	stdErrWinRate := stdWinRate / math.Sqrt(float64(nTrials))
	fmt.Println("Avg win rate =", avgWinRate, "std =", stdWinRate, "stderr =", stdErrWinRate)
}

func runGame(strategy agent.Strategy, verbose bool) (won bool) {

	deck := deck.NewDeck()
	deck.Shuffle()

	game := game.NewGame(deck)

	agent,err := agent.NewAgent(game, strategy)
	if err != nil {
		panic(err)
	}
	
	if verbose { game.Display(true) }
	
	var turnsWithoutMove int
	for turnsWithoutMove < max(len(game.Avail) + len(game.Deck), 10) {
		movedCard := agent.Act(verbose)
		if verbose { game.Display(true) }

		if movedCard {
			turnsWithoutMove = 0
		} else {
			turnsWithoutMove++
		}
	}

	if verbose {
		fmt.Println("Game is over! Final state:")
		game.Display(false)
	}
	return game.IsWon()
}