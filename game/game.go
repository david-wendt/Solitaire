package game

import (
	"solitaire/deck"
	"fmt"
)

const nSuits = deck.NSuits
const nStacks = 7
const stackCap = deck.SuitSize + nStacks - 1
type Card struct {
	id deck.Card 
	hidden bool
}

type Game struct {
	SuitStacks [nSuits]int
	PlayStacks [nStacks][]Card
	Deck []deck.Card
	Avail []deck.Card
}

func NewGame(d deck.Deck) *Game {
	game := Game{}

	cardIdx := 0

	// Initialize stacks
	for i := 0; i < nStacks; i++ {
		game.PlayStacks[i] = make([]Card, 0, stackCap)
		for j := 0; j <= i; j++ {
			var hidden bool = j < i
			game.PlayStacks[i] = append(game.PlayStacks[i], Card{id:d[cardIdx], hidden: hidden})
			cardIdx++
		}
	}

	// Initialize deck
	game.Deck = d[cardIdx:]
	game.Avail = make([]deck.Card, 0, len(game.Deck))

	return &game
}

func (game *Game) Display(hidden bool) {
	// Suit stacks
	suitStackString := " "
	for suit,stackSize := range game.SuitStacks {
		if game.SuitStacks[suit] == 0 {
			suitStackString += "_"
		} else {
			suitStackString += deck.Ranks[stackSize-1]
		}
		suitStackString += deck.Suits[suit] + "   "
	}
	fmt.Println(suitStackString)

	// Play stacks
	maxDepth := 0
	for _,stack := range game.PlayStacks {
		maxDepth = max(maxDepth, len(stack))
	}

	colIdString := ""
	for i := range game.PlayStacks {
		colIdString += fmt.Sprintf("  %v   ", i)
	}
	fmt.Println(colIdString)

	playStackString := ""
	for depth := range maxDepth {
		for _,stack := range game.PlayStacks {
			if len(stack) <= depth {
				playStackString += "      "
				continue 
			}

			card := stack[depth]
			if card.id.Rank != deck.Ten ||
				(card.hidden && hidden) {
				playStackString += " "
			} 

			if !card.hidden {
				// Card not hidden
				playStackString += fmt.Sprintf(" %v  ", card.id)
			} else {
				if hidden {
					playStackString += "[__] "
				} else {
					playStackString += fmt.Sprintf("[%v] ", card.id)
				}
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

func (game *Game) peekStack(stack int) (*Card,error) {
	stack_ := game.PlayStacks[stack]
	if len(stack_) == 0 {
		return nil,MoveError("Peeking an empty stack!")
	}
	return &stack_[len(stack_) - 1],nil
}

func (game *Game) stackColor(stack int) (byte,error) {
	card,err := game.peekStack(stack)
	if err != nil {
		return 0,err
	}
	return card.id.Color(),nil
}

func (game *Game) checkReveal(stack int) error {
	card,err := game.peekStack(stack)
	if err != nil { return err }

	if len(game.PlayStacks[stack]) > 0 && card.hidden {
		card.hidden = false
	}
	return nil 
}

// Peek the card from top stack for suit `suit`
func (game *Game) peekSuit(suit int) (Card,error) {
	stackSize := game.SuitStacks[suit]
	if stackSize == 0 {
		return Card{}, MoveError(fmt.Sprintf("Trying to pop from empty suit stack %v!", suit))
	}

	return Card{id:deck.NewCard(stackSize-1, suit)},nil
}

// Pop the card from top stack for suit `suit`
func (game *Game) popSuit(suit int) (Card,error) {
	card,err := game.peekSuit(suit)
	if err != nil {
		return Card{},err
	}

	game.SuitStacks[suit]--
	return card,nil
}

// Pop the card from stack `src`
func (game *Game) popStack(src int) (Card,error) {
	card,err := game.peekStack(src)
	if err != nil {
		return Card{},err
	}
	
	if card.hidden {
		return Card{},MoveError(fmt.Sprintf("Trying to pop hidden card %v from stack %v!", card, src))
	}

	stack := game.PlayStacks[src]
	game.PlayStacks[src] = stack[:len(stack)-1]
	game.checkReveal(src)
	return *card,nil
}

func (game *Game) peekAvail() (deck.Card,error) {
	if len(game.Avail) == 0 {
		return deck.Card{},MoveError("Trying to peek an empty Avail!")
	}
	return game.Avail[len(game.Avail) - 1],nil
}

// Pop one card from game.Avail
func (game *Game) popAvail() (Card,error) {
	card,err := game.peekAvail()
	if err != nil {
		return Card{},err
	}

	game.Avail = game.Avail[:len(game.Avail)-1]
	return Card{id:card},nil
}

// Push the card `card` to its suit stack
func (game *Game) pushSuit(card Card) error {
	suit := card.id.Suit
	if int(card.id.Rank) != game.SuitStacks[suit] {
		return MoveError(fmt.Sprintf("Invalid attempt to push card %v on %v stack of size %v!",
			card.id, suit, game.SuitStacks[suit]))
	}

	game.SuitStacks[suit]++
	return nil
}

func (game *Game) invalidPushStack(card deck.Card, dst int) (bool,error) {
	validKingMove := card.Rank == deck.King && len(game.PlayStacks[dst]) == 0
	if validKingMove {
		return false,nil
	}
	dstCard,err := game.peekStack(dst)
	if err != nil { 
		return false,err 
	}

	return ((card.Rank != dstCard.id.Rank - 1) || card.Color() == dstCard.id.Color()),nil
}

// Push the card `card` to stack `dst`
func (game *Game) pushStack(card Card, dst int) error {
	isInvalid,err := game.invalidPushStack(card.id, dst)
	if err != nil {
		return err
	}

	if isInvalid {
		dstCard,dstErr := game.peekStack(dst)
		if dstErr != nil {
			return dstErr
		}

		return MoveError(fmt.Sprintf("Invalid attempt to append card %v to stack %v with final card %v!",
			card.id, dst, dstCard.id))
	}

	game.PlayStacks[dst] = append(game.PlayStacks[dst], card)
	return nil 
}

// Move `ncards` from stack `src` to stack `dst`
func (game *Game) Move(src int, dst int, ncards int) error {
	srcStack := game.PlayStacks[src]
	
	nVisibleCards := 0
	for i := len(srcStack) - 1; i >= 0; i-- {
		if srcStack[i].hidden {
			break
		} else {
			nVisibleCards++
		}
	}

	if ncards == 0 {
		ncards = nVisibleCards
	} else if (nVisibleCards < ncards) {
		return MoveError(fmt.Sprintf(
			"Invalid move! %v cards from src %v with only %v visible cards!", ncards, src, nVisibleCards))
	}

	cards := srcStack[len(srcStack) - ncards:]

	// Need to check this ahead of time so that error is not thrown only on push,
	// leaving pop complete and game state invalid
	isInvalid,err := game.invalidPushStack(cards[0].id, dst)
	if err != nil {
		return err 
	}

	if isInvalid {
		return MoveError(fmt.Sprintf("Invalid move! %v cards from src %v to dst %v.", ncards, src, dst))
	}

	for _,card := range cards {
		game.popStack(src)
		game.pushStack(card, dst)
	}

	return nil
}

// Move one card from game.Avail to stack `dst`
func (game *Game) MoveFromAvail(dst int) error {
	card,err := game.peekAvail()
	if err != nil {
		return err
	}

	err2 := game.pushStack(Card{id:card}, dst)
	if err2 != nil {
		return err2
	}

	_,err3 := game.popAvail() // This will never fail if peekAvail does not fail
	if err3 != nil {
		panic("This should be impossible!")
	}
	return nil
}

func (game *Game) MoveAvailToTop() error {
	card,err := game.peekAvail()
	if err != nil {
		return err
	}

	err2 := game.pushSuit(Card{id:card})
	if err2 != nil {
		return err2
	}

	_,err3 := game.popAvail() // This will never fail if peekAvail does not fail
	if err3 != nil {
		panic("This should be impossible!")
	}
	return nil
}

// Move one card from stack for suit `suit` to stack `dst`
func (game *Game) MoveFromTop(suit int, dst int) error {
	card,err := game.peekSuit(suit)
	if err != nil {
		return err
	}

	err2 := game.pushStack(card, dst)
	if err2 != nil {
		return err2
	}

	_,err3 := game.popSuit(suit) // This will never fail if peekSuit does not fail
	if err3 != nil {
		panic("This should be impossible!")
	}
	return nil
}

// Move one card from stack `src` to top
func (game *Game) MoveToTop(src int) error {
	card,err := game.peekStack(src)
	if err != nil {
		return err
	}

	err2 := game.pushSuit(*card)
	if err2 != nil {
		return err2
	}

	_,err3 := game.popStack(src) // This will never fail if peekAvail does not fail
	if err3 != nil {
		panic("This should be impossible!")
	}
	return nil
}

type MoveError string 
func (err MoveError) Error() string {
	return string(err)
}