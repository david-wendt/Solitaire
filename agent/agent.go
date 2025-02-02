package agent

import (
	"fmt" 

	"solitaire/deck"
	"solitaire/game"
)

type Moves struct {
	Tableau [][2]int
	Avail []int 
	ToTop []int 
	FromTop [][2]int
}

func (moves Moves) len() int {
	return len(moves.Tableau) + len(moves.Avail) + len(moves.ToTop) + len(moves.FromTop)
}

type Agent struct {
	game *game.Game
	highCards map[deck.Card]int
	lowCards map[deck.Card]int
	emptyStack int
	strategy Strategy
}

type InitializationError string 
func (err InitializationError) Error() string {
	return string(err)
}

func NewAgent(game *game.Game, strategy Strategy) (*Agent,error) {
	agent := Agent{
		game: game, 
		highCards: make(map[deck.Card]int), 
		lowCards: make(map[deck.Card]int), 
		emptyStack: -1,
		strategy: strategy,
	}

	for i,stack := range game.HiddenStacks {
		queue := game.VisibleQueues[i]
		if len(stack) != i || len(queue) != 1 {
			return nil,InitializationError("NewAgent called on Game not in initial state!")
		}
		card := queue[0]
		agent.highCards[card] = i
		agent.lowCards[card] = i
	}

	return &agent,nil
}

func (agent *Agent) recomputeHighLowCards() {
	// Clear previous values
	for c := range agent.highCards {
		delete(agent.highCards, c)
	}
	for c := range agent.lowCards {
		delete(agent.lowCards, c)
	}
	agent.emptyStack = -1
	
	for i,queue := range agent.game.VisibleQueues {
		length := len(queue)
		if length == 0 { 
			agent.emptyStack = i // can overwrite since empty stacks are indistinguishable
		 } else {
			agent.highCards[queue[length-1]] = i
			agent.lowCards[queue[0]] = i
		 }
	}
}

func (agent *Agent) findTableauMoves() [][2]int {
	moves := make([][2]int, 0, game.NStacks - 1)

	for highCard,src := range agent.highCards {
		if highCard.Rank == deck.King {
			if agent.emptyStack != -1 && len(agent.game.HiddenStacks[src]) > 0 {
				moves = append(moves, [2]int{src, agent.emptyStack})
			}
		} else {
			for lowCard,dst := range agent.lowCards {
				if deck.CanPlace(highCard, lowCard) {
				   moves = append(moves, [2]int{src,dst})
			   }
		   }
		}
	}

	return moves 
}

func (agent *Agent) findAvailMoves() []int {
	availCard,err := agent.game.PeekAvail()
	moves := make([]int, 0, game.NStacks)
	if err == nil {
		for i := range game.NStacks {
			card,peekErr := agent.game.PeekQueue(i)
			if peekErr == nil { // nonempty stack, must be able to place card
				if deck.CanPlace(availCard, card) {
					moves = append(moves, i)
				} 
			} else { // empty stack, king is valid move
				if len(agent.game.HiddenStacks[i]) > 0 {
					panic("Empty queue should imply empty stack!")
				}
				if availCard.Rank == deck.King {
					moves = append(moves, i)
				}
			}
		}

		canPush,_ := agent.game.CanPushSuit(availCard)
		if canPush {
			moves = append(moves, -1)
		}
	}

	return moves 
}

func (agent *Agent) findMovesToTop() []int {
	moves := make([]int, 0, game.NStacks)
	for i := range game.NStacks {
		card,peekErr := agent.game.PeekQueue(i)
		if peekErr == nil {
			canPush,_ := agent.game.CanPushSuit(card)
			if canPush {
				moves = append(moves, i)
			}
		}
	}
	return moves 
}

func (agent *Agent) findMovesFromTop() [][2]int {
	moves := make([][2]int, 0, game.NStacks)
	for suitID := range deck.NSuits {
		suitCard,err := agent.game.PeekSuit(suitID)
		if err == nil && suitCard.Rank != deck.Ace {
			for i := range game.NStacks {
				card,peekErr := agent.game.PeekQueue(i)
				if peekErr == nil && deck.CanPlace(suitCard, card) {
					moves = append(moves, [2]int{suitID, i})
				}
			}
		}
	}
	return moves 
}

func (agent *Agent) findMoves() Moves {
	return Moves{
		Tableau: agent.findTableauMoves(),
		Avail: agent.findAvailMoves(),
		ToTop: agent.findMovesToTop(),
		FromTop: agent.findMovesFromTop(),
	}
}

func (agent *Agent) PrintValidMoves() {
	agent.game.Display(true)
	fmt.Println("\nAgent thinks high/low cards are:")
	fmt.Printf("High cards:%v\nLow cards:%v\n", agent.highCards, agent.lowCards)
	
	fmt.Printf("Valid moves are: %+v\n", agent.findMoves())
}

func (agent *Agent) executeMove(moves Moves, idx int) {
	if idx >= moves.len() || idx < -1 {
		panic("Tried to execute a move with invalid index!")
	}

	var err error
	// Move indexes are ordered as: Tableau, Avail, ToTop, FromTop
	l0 := len(moves.Tableau)
	l1 := len(moves.Avail)
	l2 := len(moves.ToTop)
	l3 := len(moves.FromTop) // Sadly, not Lagrange points
	if idx == -1 {
		agent.game.Flip()
	} else if idx < l0 { // Tableau move
		move := moves.Tableau[idx]
		src,dst := move[0],move[1]
		err = agent.game.Move(src, dst, 0) // WARNING! Hardcoded to move ALL cards right now.
	} else if idx - l0 < l1 { // Avail move
		dst := moves.Avail[idx - l0]
		if dst == -1 {
			err = agent.game.MoveAvailToTop()
		} else {
			err = agent.game.MoveFromAvail(dst)
		}
	} else if idx - l0 - l1 < l2 { // move ToTop
		src := moves.ToTop[idx - l0 - l1]
		err = agent.game.MoveToTop(src)
	} else if idx - l0 - l1 - l2 < l3 { // move FromTop
		move := moves.FromTop[idx - l0 - l1 - l2]
		src,dst := move[0],move[1]
		err = agent.game.MoveFromTop(src,dst)
	}

	if err != nil {
		panic(err)
	}
}

func (agent *Agent) Act(verbose bool) (movedCard bool) {
	agent.recomputeHighLowCards() // Remove this if we keep high/low cards up-to-date as intended!
	moves := agent.findMoves()
	var moveID int = -1
	if moves.len() > 0 {
		moveID = agent.strategy.choose(&moves)
	}

	if verbose {
		fmt.Printf("Valid moves found: %+v\n", moves)
		fmt.Printf("Executing move with index %v\n", moveID)
	}
	
	agent.executeMove(moves, moveID)
	movedCard = moveID != -1
	return movedCard
}
