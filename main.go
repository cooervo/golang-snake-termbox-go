package main

import (
	tb "github.com/nsf/termbox-go"
	"time"
)

var curCol = 0
var curRune = 0
var backbuf []tb.Cell
var bbw, bbh int

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

type attrFunc func(int) (rune, tb.Attribute, tb.Attribute)

func updateAndDrawButtons(current *int, x, y int, mx, my int, n int, attrf attrFunc) {
	/*
	lx, ly := x, y
	for i := 0; i < n; i++ {
		if lx <= mx && mx <= lx+3 && ly <= my && my <= ly+1 {
			*current = i
		}
		r, fg, bg := attrf(i)
		tb.SetCell(lx+0, ly+0, r, fg, bg)
		tb.SetCell(lx+1, ly+0, r, fg, bg)
		tb.SetCell(lx+2, ly+0, r, fg, bg)
		tb.SetCell(lx+3, ly+0, r, fg, bg)
		tb.SetCell(lx+0, ly+1, r, fg, bg)
		tb.SetCell(lx+1, ly+1, r, fg, bg)
		tb.SetCell(lx+2, ly+1, r, fg, bg)
		tb.SetCell(lx+3, ly+1, r, fg, bg)
		lx += 4
	}
	lx, ly = x, y
	for i := 0; i < n; i++ {
		if *current == i {
			fg := tb.ColorRed | tb.AttrBold
			bg := tb.ColorDefault
			tb.SetCell(lx+0, ly+2, '^', fg, bg)
			tb.SetCell(lx+1, ly+2, '^', fg, bg)
			tb.SetCell(lx+2, ly+2, '^', fg, bg)
			tb.SetCell(lx+3, ly+2, '^', fg, bg)
		}
		lx += 4
	}
	*/
}

var x, y = 10, 10
var key tb.Key = tb.KeyArrowRight

func redraw(mx, my int) {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)
	/*
	if mx != -1 && my != -1 {
		backbuf[bbw*my+mx] = tb.Cell{Ch: runes[curRune], Fg: colors[curCol]}
	}
	copy(tb.CellBuffer(), backbuf)
	_, h := tb.Size()
	updateAndDrawButtons(&curRune, 0, 0, mx, my, len(runes), func(i int) (rune, tb.Attribute, tb.Attribute) {
		return runes[i], tb.ColorDefault, tb.ColorDefault
	})
	updateAndDrawButtons(&curCol, 0, h-3, mx, my, len(colors), func(i int) (rune, tb.Attribute, tb.Attribute) {
		return ' ', tb.ColorDefault, colors[i]
	})
	*/
	if key == tb.KeyArrowUp {
		y--
	} else if key == tb.KeyArrowDown {
		y++
	} else if key == tb.KeyArrowLeft {
		x--
	} else if key == tb.KeyArrowRight {
		x++
	}
	tb.SetCell(x, y, runes[4], colors[5], colors[7])

	tb.Flush()
}

func reallocBackBuffer(w, h int) {
	bbw, bbh = w, h
	backbuf = make([]tb.Cell, w*h)
}

func eventListener(chEvent chan<- tb.Event) {
	chEvent <- tb.PollEvent()
}

func timeout(chTimeout chan<- string) {
	time.Sleep(time.Second * 1)
	chTimeout <- "timeout"
}

func main() {
	err := tb.Init()
	if err != nil {
		panic(err)
	}
	defer tb.Close()
	tb.SetInputMode(tb.InputEsc | tb.InputMouse)
	reallocBackBuffer(tb.Size())
	redraw(-1, -1)

mainloop:
	for range time.Tick(time.Microsecond * 500) {
		mx, my := -1, -1

		var chTimeout = make(chan string)
		var chEvent = make(chan tb.Event)

		go timeout(chTimeout)
		go eventListener(chEvent)
		print(" SEL")

		select {
		case ev := <-chEvent:
			print(" SWI ")
			switch ev.Type {
			case tb.EventKey:
				if ev.Key == tb.KeyEsc {
					break mainloop
				} else if ev.Key == tb.KeyArrowUp || ev.Key == tb.KeyArrowDown ||
					ev.Key == tb.KeyArrowRight || ev.Key == tb.KeyArrowLeft {
					key = ev.Key
				}
			case tb.EventResize:
				reallocBackBuffer(ev.Width, ev.Height)
			default:
				tb.Interrupt()
			}

			close(chEvent)
		case <-chTimeout:
			//print(" INTERRUPT ")
			//tb.Interrupt()
			close(chTimeout)
		}

		redraw(mx, my)
	}
}
