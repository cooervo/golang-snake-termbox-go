package main

import (
	tb "github.com/nsf/termbox-go"
	"time"
	"fmt"
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

type attrFunc func(int) (rune, tb.Attribute, tb.Attribute)

var x, y = 10, 10
var key tb.Key = tb.KeyArrowRight
var initDraw = true

func redraw(mx, my int) {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)

	//===
	var newHead Pos
	if key == tb.KeyArrowUp {
		newHead = Pos{x, y - 1}
	} else if key == tb.KeyArrowRight {
		newHead = Pos{x + 1, y}
	}

	// remove last element
	snake = snake[:len(snake)-1]
	var arr2 = []Pos{newHead}
	snake = append(arr2, snake...)

	//=== Change direction
	if key == tb.KeyArrowUp {
		y--
	} else if key == tb.KeyArrowRight {
		x++
	}

	// DRAW SNAKE
	for i, v := range snake {

		// HEAD
		if i == 0 {
			tb.SetCell(v.x, v.y, runes[0], tb.ColorGreen, tb.ColorGreen)
		}

		// BODY
		if i > 0 {
			tb.SetCell(v.x, v.y, runes[0], tb.ColorBlue, tb.ColorBlue)
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
	time.Sleep(time.Millisecond * 1000)
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
