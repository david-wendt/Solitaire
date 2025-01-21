package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"

	"solitaire/deck"
	"solitaire/game"
)

const nStacks int = game.NStacks

// var testing bool = true  
var testing bool = false 


func test() {
	fmt.Println("Hello world!")
}

func main() {
	if testing {
		test()
		return
	}

	deck := deck.NewDeck()
	// fmt.Println(deck)
	deck.Shuffle()
	// fmt.Println(deck)

	game := game.NewGame(deck)
	// fmt.Printf("%+v\n", game)

	game.Display(false)
	
	for {
		fmt.Print("Enter a move: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n') // Reads until newline
		move := strings.ToLower(strings.TrimSpace(input))

		if move == "h" || move == "help" {
			fmt.Println("`mv` to move a card, `f` to flip 3 cards from deck.")
			fmt.Println("With `mv`,")
			fmt.Println("\tfollow with `a t` to move from avail to top")
			fmt.Println("\tfollow with `a <dst>` to move from avail to stack <dst>")
			fmt.Println("\tfollow with `<src> <dst> <n>` to move <n> cards from stack <src> to stack <dst>")
			fmt.Println("\t\tOmit <n> to move all visible cards")
			fmt.Println("\tfollow with `<src> t` to move from stack <src> to top")
			fmt.Println("\tfollow with `t <src> <dst>` to move from top stack <src> to stack <dst>")
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
			} else if fields[0] == "a" {
				dst,ok := strconv.Atoi(fields[1])
				if ok != nil {
					fmt.Printf("Must follow `a` with int. strconv.Atoi error: %v", ok)
				} else {
					err = game.MoveFromAvail(dst)
				}
			} else if fields[0] == "t"  && len(fields) > 2{
				src,ok1 := strconv.Atoi(fields[1])
				dst,ok2 := strconv.Atoi(fields[2])
				if ok1 == nil && ok2 == nil {
					err = game.MoveFromTop(src, dst)
				} else {
					fmt.Printf("strconv.Atoi error(s): %v, %v", ok1, ok2)
				}
			} else if len(fields) > 1 && fields[1] == "t" {
				src,ok := strconv.Atoi(fields[0])
				if ok == nil {
					err = game.MoveToTop(src)
				} else {
					fmt.Printf("strconv.Atoi error: %v", ok)
				}
			} else if len(fields) > 1 {
				src,ok1 := strconv.Atoi(fields[0])
				dst,ok2 := strconv.Atoi(fields[1])
				n := 0
				var ok3 error = nil
				if len(fields) > 2 {
					n,ok3 = strconv.Atoi(fields[2])
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