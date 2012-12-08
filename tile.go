package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gcapell/dlx"
)

type board struct {
	tile
	index map[square]int // position of each square
}

type tile struct {
	squares []square
}

func (t tile) bounds() (maxX, maxY int) {
	maxX, maxY = t.squares[0].x, t.squares[0].y
	for _, s := range t.squares {
		if s.x > maxX {
			maxX = s.x
		}
		if s.y > maxY {
			maxY = s.y
		}
	}
	return maxX + 1, maxY + 1
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
	maxX, maxY := b.bounds()

	for tileIndex, t := range tiles {
		rows := 0
		for _, t2 := range permute(t) {
			for x := 0; x < maxX; x++ {
				for y := 0; y < maxY; y++ {
					t3 := t2.translate(x, y)
					if positions, ok := b.positions(t3); ok {
						positions = append(positions, len(b.squares)+tileIndex)
						rows++
						d.AddRow(positions)
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
func permute(orig tile) []tile {
	reply := []tile{orig}

	f := func(t tile, fn func(square) square) tile {
		s2 := make([]square, len(t.squares))
		copy(s2, t.squares)
		for i, s := range s2 {
			s2[i] = fn(s)
		}
		normalise(s2)
		t2 := tile{s2}
		if !contains(reply, t2) {
			reply = append(reply, t2)
		}
		return t2
	}

	r := f(orig, ror)
	r = f(r, ror)
	r = f(r, ror)

	flipped := f(orig, flip)
	r = f(flipped, ror)
	r = f(r, ror)
	r = f(r, ror)

	return reply
}

func contains(bag []tile, t tile) bool {
	for _, t2 := range bag {
		if matches(t2, t) {
			return true
		}
	}
	return false
}

func matches(t, t2 tile) bool {
	bag := make(map[square]bool)
	for _, s := range t.squares {
		bag[s] = true
	}
	for _, s := range t2.squares {
		if !bag[s] {
			return false
		}
	}
	return true
}

func normalise(ss []square) {
	minX, minY := ss[0].x, ss[0].y

	for _, s := range ss {
		if s.x < minX {
			minX = s.x
		}
		if s.y < minY {
			minY = s.y
		}
	}
	if minX == 0 && minY == 0 {
		return
	}
	for i, s := range ss {
		ss[i] = square{s.x - minX, s.y - minY}
	}
}

// Rotate tile quarter-turn clockwise
func ror(s square) square {
	return square{s.y, -s.x}
}

func flip(s square) square {
	return square{-s.x, s.y}
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

func (b *board) positions(t tile) ([]int, bool) {
	reply := make([]int, len(t.squares))
	var ok bool
	for i, s := range t.squares {
		if reply[i], ok = b.index[s]; !ok {
			return reply, false
		}
	}
	return reply, true
}

// Return string representation of solution
func (b *board) text(nTiles int, dlxResult [][]int) string {

	maxX, maxY := b.bounds()

	letters := make([][]byte, maxY)
	for j := 0; j < maxY; j++ {
		letters[j] = make([]byte, maxX)
		for i := range letters[j] {
			letters[j][i] = ' '
		}
	}

	for _, row := range dlxResult {
		tile := row[len(row)-1] - len(b.squares)
		for _, val := range row[:len(row)-1] {
			s := b.squares[val]
			letters[s.y][s.x] = byte('a' + tile)
		}
	}
	result := ""
	for _, line := range letters {
		result += string(line) + "\n"
	}

	return result
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

	}
	return squares, nil
}
