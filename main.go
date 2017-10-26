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

func update_and_redraw_all(mx, my int) {
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
var ev tb.Event
func listenEventOrTimeout() {
	print("timeout START")
	time.Sleep(time.Second * 2)
	print("timeout FINISHED")

	 tb.Interrupt()
	print("tbInterrupt()")
}
func main() {
	err := tb.Init()
	if err != nil {
		panic(err)
	}
	defer tb.Close()
	tb.SetInputMode(tb.InputEsc | tb.InputMouse)
	reallocBackBuffer(tb.Size())
	update_and_redraw_all(-1, -1)

mainloop:
	for {
		mx, my := -1, -1
		print(0)
		ev := tb.PollEvent()
		print(1)
		listenEventOrTimeout()
		print(2)
		switch ev.Type {
		case tb.EventKey:
			print("case EVENTKEY")
			if ev.Key == tb.KeyEsc {
				break mainloop
			} else if ev.Key == tb.KeyArrowUp || ev.Key == tb.KeyArrowDown ||
				ev.Key == tb.KeyArrowRight || ev.Key == tb.KeyArrowLeft {
				key = ev.Key
			}
		case tb.EventResize:
			print("case EVENT RIsE")
			reallocBackBuffer(ev.Width, ev.Height)
		case tb.EventInterrupt:
			print("case EVENT INTERRUPT")

		default:
			print("DEF")
			tb.Interrupt()
		}

		print("under switch")
		update_and_redraw_all(mx, my)
	}
}
