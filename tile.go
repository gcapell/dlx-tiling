package main

import (
	"errors"
	"fmt"
	"github.com/gcapell/dlx"
	"log"
	"strings"
)

type board struct {
	nSquares int
	squares  []square
	nTiles   int
}

type tile struct {
	index   int
	squares []square
}

type square struct {
	x, y int
}

func main() {
	tiles := asciiToTiles(tilesASCII)
	b := parseBoard(boardASCII)
	err, result := b.solve(tiles)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

func (b *board) solve(tiles []tile) (error, string) {
	if problem := b.sanityCheck(tiles); problem != nil {
		return problem, ""
	}

	b.nTiles = len(tiles)
	d := dlx.New(b.nTiles + b.nSquares)
	for _, t := range tiles {
		for _, t2 := range permute(t) {
			for _, p := range b.squares {
				t3 := t2.translate(p)
				if b.contains(t3) {
					d.AddRow(b.dlxRow(t3))
				}
			}
		}
	}
	result := d.Search()
	if result == nil {
		return errors.New("No solution found"), ""
	}
	return nil, b.text(result)
}

// Permute returns all the distinct permutations (rotate/flip/...) of a tile.
func permute(t tile) []tile {
	return nil // fixme
}

func (t tile) translate(offset square) tile {
	return t // fixme
}

func (b *board) contains(t tile) bool {
	return true // fixme
}

func (b *board) dlxRow(t tile) []int {
	return nil // fixme
}

// Return string representation of solution
func (b *board) text(dlxResult [][]int) string {
	return "cannot represent solution" // fixme
}

func (b *board) sanityCheck(tiles []tile) error {
	s := 0 // Total squares covered by all tiles
	for _, t := range tiles {
		s += len(t.squares)
	}
	if s != b.nSquares {
		return fmt.Errorf("tiles cover %d squares, board has %d squares",
			s, b.nSquares)
	}
	return nil
}

// Parse ascii drawing of board
func parseBoard(s string) board {
	return board{} // fixme
}

// Parse (blank-line-separated) ascii drawings of tiles
func asciiToTiles(s string) []tile {
	chunks := strings.Split(s, "\n\n")
	log.Printf("%d chunks, :%#v\n", len(chunks), chunks)
	tiles := make([]tile, 0, len(chunks))
	index := 0
	for _, c := range chunks {
		squares, err := asciiToSquares(c)
		if err != nil {
			log.Println(err)
			continue
		}
		tiles = append(tiles, tile{index, squares})
		index++
	}
	return tiles
}

// Parse ascii drawing of tile
func asciiToSquares(s string) ([]square, error) {
	squares := make([]square, 0, 6)
	chunks := strings.Split(s, "\n")
	row := 0
	for _, chunk := range chunks {
		if len(chunk) == 0 {
			continue
		}
		for col, c := range chunk {
			switch c {
			case 'x':
				squares = append(squares, square{col, row})
			case ' ':
				continue
			default:
				return squares, fmt.Errorf("unrecognised char %c", c)
			}
		}
		row++
		log.Println()
	}
	log.Printf("asciiToSquares: %#v->%v\n", chunks, squares)
	return squares, nil // FIXME
}
