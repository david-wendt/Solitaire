package agent

import "fmt"

type Moves struct {
	Tableau [][2]int
	Avail []int 
	ToTop []int 
	FromTop [][2]int
}

func (moves Moves) len() int {
	return len(moves.Tableau) + len(moves.Avail) + len(moves.ToTop) + len(moves.FromTop)
}

type MoveID interface { int | [2]int }
func index[T MoveID](slice []T, elt T) int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == elt {
			return i
		}
	}
	return -1
}

func (moves Moves) index(otherMoves Moves) int {
	if otherMoves.len() != 1 { panic(fmt.Sprintf("Trying to index Moves with argument with length %v != 1!", otherMoves.len())) }

	totalIdx := 0
	if len(otherMoves.Tableau) == 1 {
		otherMove := otherMoves.Tableau[0]
		idx := index(moves.Tableau, otherMove)
		if idx == -1 {
			return -1
		} else {
			return totalIdx + idx 
		}
	} else {
		totalIdx += len(moves.Tableau)
	}
	
	if len(otherMoves.Avail) == 1 {
		otherMove := otherMoves.Avail[0]
		idx := index(moves.Avail, otherMove)
		if idx == -1 {
			return -1
		} else {
			return totalIdx + idx 
		}
	} else {
		totalIdx += len(moves.Avail)
	}
	
	if len(otherMoves.ToTop) == 1 {
		otherMove := otherMoves.ToTop[0]
		idx := index(moves.ToTop, otherMove)
		if idx == -1 {
			return -1
		} else {
			return totalIdx + idx 
		}
	} else {
		totalIdx += len(moves.ToTop)
	}
	
	if len(otherMoves.FromTop) == 1 {
		otherMove := otherMoves.FromTop[0]
		idx := index(moves.FromTop, otherMove)
		if idx == -1 {
			return -1
		} else {
			return totalIdx + idx 
		}
	}

	return -1
}