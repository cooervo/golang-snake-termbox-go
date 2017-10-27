package main

import (
	tb "github.com/nsf/termbox-go"
	"time"
	"fmt"
	"math/rand"
)

var backbuf []tb.Cell

var corners = []rune{'O', '░', '▒', '▓', '█'}
var runes = []rune{' ', '░', '▒', '▓', '█'}
var colors = []tb.Attribute{
	tb.ColorBlack,
	tb.ColorRed,
	tb.ColorGreen,
	tb.ColorYellow,
	tb.ColorBlue,
	tb.ColorMagenta,
	tb.ColorCyan,
	tb.ColorWhite,
}

type Pos struct {
	x, y int
}

func (p Pos) String() string {
	return fmt.Sprintf("%p", p)
}

var snake = []Pos{Pos{x, y}, Pos{x - 1, y}, Pos{x - 2, y}}
var applePos = Pos{5, 5}

type attrFunc func(int) (rune, tb.Attribute, tb.Attribute)

var x, y = 10, 10
var key tb.Key = tb.KeyArrowRight
var initDraw = true

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func redraw(mx, my int) {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)

	//===
	var newHead Pos
	if key == tb.KeyArrowUp {
		newHead = Pos{x, y - 1}
	} else if key == tb.KeyArrowDown {
		newHead = Pos{x, y + 1}
	} else if key == tb.KeyArrowRight {
		newHead = Pos{x + 1, y}
	} else if key == tb.KeyArrowLeft {
		newHead = Pos{x - 1, y}
	}

	// === APPLe
	tb.SetCell(applePos.x, applePos.y, '🍎', tb.ColorRed, tb.ColorBlack)

	if newHead.x == applePos.x && newHead.y == applePos.y {
		applePos = Pos{random(0, 20), random(0, 20)}

		var newTail Pos
		if key == tb.KeyArrowUp {
			newTail = Pos{x, y - 1}
		} else if key == tb.KeyArrowDown {
			newTail = Pos{x, y + 1}
		} else if key == tb.KeyArrowRight {
			newTail = Pos{x + 1, y}
		} else if key == tb.KeyArrowLeft {
			newTail = Pos{x - 1, y}
		}

		snake = append(snake, newTail)
	}

	// remove last element
	snake = snake[:len(snake)-1]
	// create new slice with newHead as first element
	snake = append([]Pos{newHead}, snake...)

	//=== Change direction
	if key == tb.KeyArrowUp {
		y--
	} else if key == tb.KeyArrowDown {
		y++
	} else if key == tb.KeyArrowRight {
		x++
	} else if key == tb.KeyArrowLeft {
		x--
	}

	// DRAW SNAKE
	for i, v := range snake {
		// HEAD
		if i == 0 {
			tb.SetCell(v.x, v.y, '▣', tb.ColorGreen, tb.ColorBlack)
		}

		// BODY
		if i > 0 {
			tb.SetCell(v.x, v.y, '▣', tb.ColorGreen, tb.ColorBlack)
		}

	}

	tb.Flush()
}

var width int
var height int

func reallocBackBuffer(w, h int) {
	width = w
	height = h
	backbuf = make([]tb.Cell, w*h)
}

func eventListener(chEvent chan<- tb.Event) {
	chEvent <- tb.PollEvent()
}

func timeout(chTimeout chan<- string) {
	time.Sleep(time.Millisecond * 300)
	chTimeout <- "timeout"
}

func main() {
	err := tb.Init()
	if err != nil {
		panic(err)
	}
	defer tb.Close()
	reallocBackBuffer(tb.Size())

snakeLoop:
	for {
		mx, my := -1, -1

		var chTimeout = make(chan string)
		var chEvent = make(chan tb.Event)

		go timeout(chTimeout)
		go eventListener(chEvent)

		select {
		case ev := <-chEvent:
			if ev.Key == tb.KeyEsc {
				break snakeLoop

			} else if ev.Key == tb.KeyArrowUp || ev.Key == tb.KeyArrowDown ||
				ev.Key == tb.KeyArrowRight || ev.Key == tb.KeyArrowLeft {
				key = ev.Key

			} else if ev.Type == tb.EventResize {
				reallocBackBuffer(ev.Width, ev.Height)
			}

			close(chEvent)
		case <-chTimeout:
			close(chTimeout)
		}

		redraw(mx, my)
	}
}
