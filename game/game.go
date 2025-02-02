package game

import (
	"solitaire/deck"
	"fmt"
	"strings"
)

const nSuits = deck.NSuits
const NStacks = 7
// const stackCap = deck.SuitSize + NStacks - 1

type Game struct {
	SuitStacks [nSuits]int
	HiddenStacks [NStacks][]deck.Card // End of slice is top of stack
	VisibleQueues [NStacks][]deck.Card // End of slice is front of queue (i.e. bottom of "stack")
	Deck []deck.Card
	Avail []deck.Card
}

func NewGame(d deck.Deck) *Game {
	game := Game{}

	cardIdx := 0

	// Initialize stacks
	for i := 0; i < NStacks; i++ {
		game.HiddenStacks[i] = make([]deck.Card, 0, i)
		game.VisibleQueues[i] = make([]deck.Card, 0, NStacks)
		for j := 0; j <= i; j++ {
			var hidden bool = j < i
			if hidden {
				game.HiddenStacks[i] = append(game.HiddenStacks[i], d[cardIdx])
			} else {
				game.VisibleQueues[i] = append(game.VisibleQueues[i], d[cardIdx])
			}
			
			cardIdx++
		}

		if len(game.HiddenStacks[i]) != i || len(game.VisibleQueues[i]) != 1 {
			panic(fmt.Sprintf("Broken implementation! Stack %v has %v hidden cards and %v visible cards",
				i, len(game.HiddenStacks[i]), len(game.VisibleQueues[i])))
		}
	}

	// Initialize deck
	game.Deck = d[cardIdx:]
	game.Avail = make([]deck.Card, 0, len(game.Deck))

	return &game
}

func (game *Game) Display(hidden bool) {
	// Suit stacks
	suitStackString := "  "
	for suit,stackSize := range game.SuitStacks {
		if game.SuitStacks[suit] == 0 {
			suitStackString += "_"
		} else {
			suitStackString += deck.Ranks[stackSize-1]
		}
		suitStackString += deck.Suits[suit] + "    "
	}
	fmt.Println(suitStackString)

	// Play stacks
	maxDepth := 0
	for i := range NStacks {
		stackLen := len(game.HiddenStacks[i])
		queueLen := len(game.VisibleQueues[i])
		maxDepth = max(maxDepth, stackLen + queueLen)
	}
	fmt.Println(strings.Repeat("-", 43))

	colIdString := ""
	for i := range NStacks {
		colIdString += fmt.Sprintf("  %v   ", i)
	}
	fmt.Println(colIdString)

	playStackString := ""
	for depth := range maxDepth {
		for i := range NStacks {
			hiddenStack := game.HiddenStacks[i]
			visibleQueue := game.VisibleQueues[i]

			if depth < len(hiddenStack) {
				if hidden {
					playStackString += " [__] "
				} else {
					card := hiddenStack[depth]
					if card.Rank != deck.Ten {
						playStackString += " "
					}
					playStackString += fmt.Sprintf("[%v] ", card)
				}
			} else if depth < len(hiddenStack) + len(visibleQueue) {
				card := visibleQueue[len(hiddenStack) + len(visibleQueue) - depth - 1]

				if card.Rank != deck.Ten {
					playStackString += " "
				}
				playStackString += fmt.Sprintf(" %v  ", card)
			} else {
				playStackString += "      "
				continue 
			}
		}
		playStackString += "\n"
	}
	fmt.Print(playStackString)

	fmt.Printf("Deck: [[%v]]\n", len(game.Deck))
	availString := "["
	if len(game.Avail) == 0 {
		availString += "]"
	} else {
		for i := len(game.Avail) - 1; i >= 0; i-- {
			availString += game.Avail[i].String() + "]"
		}
	}
	fmt.Printf("Avail: %v\n", availString)
}

func (game *Game) Flip() {
	// Flip 3 cards from the game.Deck to the game.Avail

	if len(game.Deck) == 0 {
		// Deck is out, swap Deck and Avail
		game.Deck = game.Avail
		game.Avail = make([]deck.Card, 0, len(game.Deck))
	}

	l := min(len(game.Deck), 3)
	game.Avail = append(game.Avail, game.Deck[:l]...)
	game.Deck = game.Deck[l:]
}

func (game *Game) PeekQueue(queueID int) (deck.Card, error) {
	queue := game.VisibleQueues[queueID]
	if len(queue) == 0 {
		return deck.Card{},MoveError(fmt.Sprintf("Failed attempt to peek VisQueue %v with no visible cards!", queueID))
	}
	return queue[0],nil
}

// func (game *Game) stackColor(stackID int) (byte,error) {
// 	card,err := game.PeekQueue(stackID)
// 	if err != nil {
// 		return 0,err
// 	}
// 	return card.Color(),nil
// }

// Peek the card from top stack for suit `suit`
func (game *Game) PeekSuit(suit int) (deck.Card,error) {
	stackSize := game.SuitStacks[suit]
	if stackSize == 0 {
		return deck.Card{}, MoveError(fmt.Sprintf("Trying to pop from empty suit stack %v!", suit))
	}

	return deck.NewCard(stackSize-1, suit),nil
}

// Pop the card from top stack for suit `suit`
func (game *Game) popSuit(suit int) (deck.Card,error) {
	card,err := game.PeekSuit(suit)
	if err != nil {
		return deck.Card{},err
	}

	game.SuitStacks[suit]--
	return card,nil
}

func (game *Game) PeekAvail() (deck.Card,error) {
	if len(game.Avail) == 0 {
		return deck.Card{},MoveError("Trying to peek an empty Avail!")
	}
	return game.Avail[len(game.Avail) - 1],nil
}

// Pop one card from game.Avail
func (game *Game) popAvail() (deck.Card,error) {
	card,err := game.PeekAvail()
	if err != nil {
		return deck.Card{},err
	}

	game.Avail = game.Avail[:len(game.Avail)-1]
	return card,nil
}

func (game *Game) CanPushSuit(card deck.Card) (bool,deck.SuitT) {
	suit := card.Suit
	canPush := int(card.Rank) == game.SuitStacks[suit]
	return canPush,suit
}

// Push the card `card` to its suit stack
func (game *Game) pushSuit(card deck.Card) error {
	canPush,suit := game.CanPushSuit(card)
	if !canPush {
		return MoveError(fmt.Sprintf("Invalid attempt to push card %v on %v stack of size %v!",
			card, suit, game.SuitStacks[suit]))
	}

	game.SuitStacks[suit]++
	return nil
}

// Peek the last `n` cards from queue `src`
// func (game *Game) PeekQueue(src int, n int) ([]deck.Card,error) {}

// Pop the last `n` cards from queue `src`
func (game *Game) popQueue(src int, n int) ([]deck.Card,error) {
	queue := game.VisibleQueues[src]
	var emptySlice []deck.Card
	if len(queue) < n {
		return emptySlice,MoveError(fmt.Sprintf(
			"Invalid popQueue! Trying to pop %v cards from stack %v with %v visible cards!",
			n, src, len(queue),
		))
	}

	game.VisibleQueues[src] = game.VisibleQueues[src][n:]
	if len(game.VisibleQueues[src]) == 0 {
		nHidden := len(game.HiddenStacks[src])
		if nHidden > 0 {
			card := game.HiddenStacks[src][nHidden - 1]
			game.HiddenStacks[src] = game.HiddenStacks[src][:nHidden-1]
			newQueue := make([]deck.Card, 0, deck.SuitSize)
			game.VisibleQueues[src] = append(newQueue, card)
		}
	} else {
		newQueue := make([]deck.Card, 0, deck.SuitSize)
		game.VisibleQueues[src] = append(newQueue, game.VisibleQueues[src]...)
	}
	return queue[:n],nil
}

func (game *Game) validPushQueue(cards []deck.Card, dst int) (bool,error) {
	ncards := len(cards)
	if ncards == 0 {
		return false,MoveError("Trying to push 0 cards!")
	}

	card := cards[ncards - 1]
	if len(game.VisibleQueues[dst]) == 0 && len(game.HiddenStacks[dst]) != 0 {
		panic(fmt.Sprintf(
			"Bug found! VisQueue[%v] has length 0, but HiddenStack[%v] has length %v!",
			dst, dst, len(game.HiddenStacks[dst]),
		))
	}

	validKingMove := card.Rank == deck.King && len(game.VisibleQueues[dst]) == 0
	if validKingMove {
		return true,nil
	}
	
	dstCard,err := game.PeekQueue(dst)
	if err != nil { 
		return false,MoveError("Invalid move! Attempted to move non-king to empty stack.") 
	}

	return ((card.Rank == dstCard.Rank - 1) && card.Color() != dstCard.Color()),nil
}

// Push the cards `cards` to stack `dst`
func (game *Game) extendQueue(cards []deck.Card, dst int) error {
	isValid,err := game.validPushQueue(cards, dst)
	if err != nil {
		return err
	}

	if !isValid {
		dstCard,dstErr := game.PeekQueue(dst)
		if dstErr != nil {
			return dstErr
		}

		return MoveError(fmt.Sprintf("Invalid attempt to append cards %v to stack %v with final card %v!",
			cards, dst, dstCard))
	}

	game.VisibleQueues[dst] = append(cards, game.VisibleQueues[dst]...)
	return nil 
}

// Push the card `card` to stack `dst`
func (game *Game) pushQueue(card deck.Card, dst int) error {
	newQueue := make([]deck.Card, 0, deck.SuitSize)
	newQueue = append(newQueue, card)
	return game.extendQueue(newQueue, dst)
}

// Move `ncards` from stack `src` to stack `dst`
func (game *Game) Move(src int, dst int, ncards int) error {
	// srcStack := game.PlayStacks[src]
	
	nVisibleCards := len(game.VisibleQueues[src])

	if ncards == 0 {
		ncards = nVisibleCards
	} else if (nVisibleCards < ncards) {
		return MoveError(fmt.Sprintf(
			"Invalid move! %v cards from src %v with only %v visible cards!", ncards, src, nVisibleCards))
	}

	valid,err0 := game.validPushQueue(game.VisibleQueues[src][:ncards], dst)
	if err0 != nil {
		return err0
	} else if !valid {
		return MoveError(fmt.Sprintf(
			"Invalid move! Cannot move %v cards from stack %v to stack %v!",
			ncards, src, dst,
		))
	}

	cards,err1 := game.popQueue(src, ncards)
	if err1 != nil {
		return err1
	}

	err2 := game.extendQueue(cards, dst)
	if err2 != nil {
		return err2
	}

	return nil
}

// Move one card from game.Avail to stack `dst`
func (game *Game) MoveFromAvail(dst int) error {
	card,err := game.PeekAvail()
	if err != nil {
		return err
	}

	err2 := game.pushQueue(card, dst)
	if err2 != nil {
		return err2
	}

	_,err3 := game.popAvail() // This will never fail if PeekAvail does not fail
	if err3 != nil {
		panic("This should be impossible!")
	}
	return nil
}

func (game *Game) MoveAvailToTop() error {
	card,err := game.PeekAvail()
	if err != nil {
		return err
	}

	err2 := game.pushSuit(card)
	if err2 != nil {
		return err2
	}

	_,err3 := game.popAvail() // This will never fail if PeekAvail does not fail
	if err3 != nil {
		panic("This should be impossible!")
	}
	return nil
}

// Move one card from stack for suit `suit` to stack `dst`
func (game *Game) MoveFromTop(suit int, dst int) error {
	card,err := game.PeekSuit(suit)
	if err != nil {
		return err
	}

	err2 := game.pushQueue(card, dst)
	if err2 != nil {
		return err2
	}

	_,err3 := game.popSuit(suit) // This will never fail if PeekSuit does not fail
	if err3 != nil {
		panic("This should be impossible!")
	}
	return nil
}

// Move one card from stack `src` to top
func (game *Game) MoveToTop(src int) error {
	card,err := game.PeekQueue(src)
	if err != nil {
		return err
	}

	err2 := game.pushSuit(card)
	if err2 != nil {
		return err2
	}

	cards,err3 := game.popQueue(src, 1) // This will never fail if PeekAvail does not fail
	if err3 != nil || len(cards) != 1 || cards[0] != card {
		panic(fmt.Sprintf("This should be impossible! %v, %v, %v", err3, cards, card))
	}
	return nil
}

type MoveError string 
func (err MoveError) Error() string {
	return string(err)
}

func (game *Game) IsWon() bool {
	won := true 
	for _,suitStackSize := range game.SuitStacks {
		won = won && (suitStackSize == deck.SuitSize)
	}
	return won 
}