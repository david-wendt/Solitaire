package main

import (
	"fmt"
	"strconv"
	"strings"

	"solitaire/deck"
	"solitaire/game"
	"solitaire/ioutils"
)

func main() {

	deck := deck.NewDeck()
	deck.Shuffle()

	game := game.NewGame(deck)

	game.Display(true) // `true` to hide hidden cards

	for {
		move := ioutils.Input("Enter a move: ")

		if move == "h" || move == "help" {
			fmt.Println("<Enter> to flip 3 cards from deck. To move cards,")
			fmt.Println("\t`a t` to move from available (waste pile) to top (foundation)")
			fmt.Println("\t`a <dst>` to move from available (waste pile) to stack <dst> (in tableau)")
			fmt.Println("\tfollow with `<src> <dst> <n>` to move <n> cards from stack <src> to stack <dst> (in tableau)")
			fmt.Println("\t\tOmit <n> to move all visible cards")
			fmt.Println("\tfollow with `<src> t` to move from stack <src> (in tableau) to top (foundation)")
			fmt.Println("\tfollow with `t <src> <dst>` to move from top (foundation) stack <src> to stack <dst> (in tableau)")
			continue
		}

		fields := strings.Fields(move)
		// fmt.Println("Parsed command:", fields)
		if move == "f" || move == "" {
			game.Flip()
		} else { // Move
			var err error
			if move == "a t" {
				err = game.MoveAvailToTop()
			} else if fields[0] == "a" && len(fields) > 1 {
				dst, ok := strconv.Atoi(fields[1])
				if ok != nil {
					fmt.Printf("Must follow `a` with int. strconv.Atoi error: %v", ok)
				} else {
					err = game.MoveFromAvail(dst)
				}
			} else if fields[0] == "t" && len(fields) > 2 {
				src, ok1 := strconv.Atoi(fields[1])
				dst, ok2 := strconv.Atoi(fields[2])
				if ok1 == nil && ok2 == nil {
					err = game.MoveFromTop(src, dst)
				} else {
					fmt.Printf("strconv.Atoi error(s): %v, %v", ok1, ok2)
				}
			} else if len(fields) > 1 && fields[1] == "t" {
				src, ok := strconv.Atoi(fields[0])
				if ok == nil {
					err = game.MoveToTop(src)
				} else {
					fmt.Printf("strconv.Atoi error: %v", ok)
				}
			} else if len(fields) > 1 {
				src, ok1 := strconv.Atoi(fields[0])
				dst, ok2 := strconv.Atoi(fields[1])
				n := 0
				var ok3 error = nil
				if len(fields) > 2 {
					n, ok3 = strconv.Atoi(fields[2])
				}
				if ok1 == nil && ok2 == nil && ok3 == nil {
					err = game.Move(src, dst, n)
				} else {
					fmt.Println("Invalid move command ! Entered non-integer <src>/<dst>/<n>.")
				}
			} else {
				fmt.Println("Invalid move command! See `help`.")
			}

			if err != nil {
				fmt.Printf("Move error: %v\n", err)
			}
		}

		game.Display(true)
	}
}
