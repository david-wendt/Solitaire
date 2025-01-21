package agent

import (
	// "fmt"
	"math/rand"
)

type Strategy interface {
	choose(*Moves) int
}

type ProbabilisticStrategy struct{
	PFlip float32
	PTableau float32
	PAvail float32
	PToTop float32
	PFromTop float32
}

func (strat ProbabilisticStrategy) choose(moves *Moves) int {
	nTableau := len(moves.Tableau)
	nAvail := len(moves.Avail)
	nToTop := len(moves.ToTop)
	nFromTop := len(moves.FromTop)

	// mp for "masked probabilities", mp = p if n > 0 else 0
	var mpTableau, mpAvail, mpToTop, mpFromTop float32 = 0., 0., 0., 0.
	if nTableau > 0 { mpTableau = strat.PTableau }
	if nAvail > 0 { mpAvail = strat.PAvail }
	if nToTop > 0 { mpToTop = strat.PToTop }
	if nFromTop > 0 { mpFromTop = strat.PFromTop }

	pTot := strat.PFlip + mpTableau + mpAvail + mpToTop + mpFromTop
	mpFlip := strat.PFlip / pTot
	mpTableau = mpTableau / pTot
	mpAvail = mpAvail / pTot
	mpToTop = mpToTop / pTot
	// mpFromTop = mpFromTop / pTot // Unecessary, since unused, and all add to 1 now

	r := rand.Float32()
	if r < mpFlip {
		return -1
	} else if r < mpFlip + mpTableau {
		return rand.Intn(nTableau)
	} else if r < mpFlip + mpTableau + mpAvail {
		return nTableau + rand.Intn(nAvail)
	} else if r < mpFlip + mpTableau + mpAvail + mpToTop {
		return nTableau + nAvail + rand.Intn(nToTop)
	} else {
		return nTableau + nAvail + nToTop + rand.Intn(nFromTop)
	}
}