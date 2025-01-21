# Notes on Solitaire Agents

Part of this project is to estimate the fraction of Solitaire games that are winnable (or at least a lower bound on this quantity) by building agents to play solitaire.

## Agents and their win rates

First agent:
```
strategy := agent.ProbabilisticStrategy{
    PFlip: 0.001, // Warning! If this is zero, can get stuck in infinite loop
    PTableau: 0.099,
    PAvail: 0.89, // This does not do what I want, since
    PToTop: 0.1, // Avail includes AvailToTop.
    PFromTop: 0.0,
}
```
I played 3 games by hand, and saw 1 win. (Note that I stopped playing once I saw a win...) Suggests a win rate of a bit under 1/3. (Exercise for the reader: What win rate does it suggest, assuming independent Bernoulli trials?)

## TODOs
* I have completed building an agent capable of playing Solitaire to completion and winning! 
    - Next step: make the agent detect the game end via a complete cycle of waste flips without any card moves (this may catch agents skipping possible moves, but that is desired if the agent will never make another move). Then, if all foundation piles are complete, count a win, otherwise a loss.