package main

import (
	tb "github.com/nsf/termbox-go"
	"time"
	"fmt"
	"math/rand"
)

var backbuf []tb.Cell

type Pos struct {
	x, y int
}

func (p Pos) String() string {
	return fmt.Sprintf("%p", p)
}

var snake = []Pos{Pos{x, y}, Pos{x - 1, y}, Pos{x - 2, y}}
var applePos = Pos{5, 5}

type attrFunc func(int) (rune, tb.Attribute, tb.Attribute)

var key tb.Key = tb.KeyArrowRight

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

var x, y = 10, 10

func redraw() {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)
	head := newHeadDirection()
	appleHandler(head)
	advanceSnake(head)
	changeDirection()
	drawSnake()
	tb.Flush()
}

func appleHandler(head Pos) {
	tb.SetCell(applePos.x, applePos.y, 'o', tb.ColorRed, tb.ColorBlack)
	if head.x == applePos.x && head.y == applePos.y {
		applePos = Pos{random(1, 20), random(1, 20)}

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
}

func newHeadDirection() Pos {
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
	return newHead
}

func advanceSnake(newHead Pos) {
	snake = snake[:len(snake)-1]
	snake = append([]Pos{newHead}, snake...)
}
func changeDirection() {
	if key == tb.KeyArrowUp {
		y--
	} else if key == tb.KeyArrowDown {
		y++
	} else if key == tb.KeyArrowRight {
		x++
	} else if key == tb.KeyArrowLeft {
		x--
	}
}

func drawSnake() {
	for i, v := range snake {
		// HEAD
		if i == 0 {
			tb.SetCell(v.x, v.y, 'O', tb.ColorGreen, tb.ColorBlack)
		}

		// BODY
		if i > 0 {
			tb.SetCell(v.x, v.y, 'O', tb.ColorGreen, tb.ColorBlack)
		}
	}
}

func reallocBackBuffer(w, h int) {
	backbuf = make([]tb.Cell, w*h)
}

func eventListener(chEvent chan<- tb.Event) {
	chEvent <- tb.PollEvent()
}

func timeout(chTimeout chan<- string) {
	time.Sleep(time.Millisecond * 200)
	chTimeout <- "timeout"
}

func main() {
	err := tb.Init()
	if err != nil {
		panic(err)
	}
	defer tb.Close()

	reallocBackBuffer(tb.Size())
	gameLoop()
}

func gameLoop() {
	for {
		var chTimeout = make(chan string)
		var chEvent = make(chan tb.Event)

		go timeout(chTimeout)
		go eventListener(chEvent)

		select {
		case ev := <-chEvent:
			// if esc cancel game
			if ev.Key == tb.KeyEsc {
				return

				// if is arrow key
			} else if ev.Key == tb.KeyArrowUp || ev.Key == tb.KeyArrowDown ||
				ev.Key == tb.KeyArrowRight || ev.Key == tb.KeyArrowLeft {
				key = ev.Key

				// if is window resize
			} else if ev.Type == tb.EventResize {
				reallocBackBuffer(ev.Width, ev.Height)
			}

			close(chEvent)
		case <-chTimeout:
			close(chTimeout)
		}

		redraw()
	}
}
