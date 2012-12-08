package main

import (
	"errors"
	"fmt"
	"github.com/gcapell/dlx"
	"log"
	"strings"
)

type board struct {
	tile
	index map[square]int	// position of each square
}

type tile struct {
	squares []square
}

func (t tile) bounds() (minX, minY, maxX, maxY int) {
	for i, s := range t.squares {
		if i == 0 {
			minX, minY, maxX, maxY = s.x, s.y, s.x, s.y
			continue
		}
		if s.x < minX {
			minX = s.x
		}
		if s.x > maxX {
			maxX = s.x
		}
		if s.y < minY {
			minY = s.y
		}
		if s.y > maxY {
			maxY = s.y
		}
	}
	return
}

type square struct {
	x, y int
}

func main() {
	tiles := asciiToTiles(tilesASCII)
	b, err := parseBoard(boardASCII)
	if err != nil {
		log.Fatal(err)
	}
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

	nTiles := len(tiles)
	d := dlx.New(nTiles + len(b.squares))
	_, maxX, _, maxY := b.bounds()

	for tileIndex, t := range tiles {
		for _, t2 := range permute(t) {
			_, tmaxX, _, tmaxY := t2.bounds()
			for x := 0; x < maxX-tmaxX; x++ {
				for y := 0; y < maxY-tmaxY; y++ {
					t3 := t2.translate(x, y)
					if b.contains(t3) {
						d.AddRow(b.dlxRow(nTiles, tileIndex, t3))
					}
				}
			}
		}
	}
	result := d.Search()
	if result == nil {
		return errors.New("No solution found"), ""
	}
	return nil, b.text(nTiles, result)
}

// Permute returns all the distinct permutations (rotate/flip/...) of a tile.
func permute(t tile) []tile {
	return nil // fixme
}

func (t tile) translate(dx, dy int) tile {
	s2 := make([]square, len(t.squares))
	copy(s2, t.squares)
	for i := range s2 {
		s2[i].x += dx
		s2[i].y += dy
	}
	return tile{s2}
}

func (b *board) contains(t tile) bool {
	for _, s := range t.squares {
		if _, ok := b.index[s]; !ok {
			return false
		}
	}
	return true // fixme
}

func (b *board) dlxRow(nTiles, tileIndex int, t tile) []int {
	return nil // fixme
}

// Return string representation of solution
func (b *board) text(nTiles int, dlxResult [][]int) string {
	return "cannot represent solution" // fixme
}

func (b *board) sanityCheck(tiles []tile) error {
	s := 0 // Total squares covered by all tiles
	for _, t := range tiles {
		s += len(t.squares)
	}
	if s != len(b.squares) {
		return fmt.Errorf("tiles cover %d squares, board has %d squares",
			s, len(b.squares))
	}
	return nil
}

// Parse ascii drawing of board
func parseBoard(s string) (*board, error) {
	squares, err := asciiToSquares(s)
	if err != nil {
		return nil, err
	}
	index := make(map[square]int)
	for i, s := range squares {
		index[s] = i
	}
	return &board{tile{squares}, index}, nil
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
		tiles = append(tiles, tile{squares})
		index++
	}
	log.Printf("Tiles: %v", tiles)
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
