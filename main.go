package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

var speed int

func main() {

	flag.IntVar(&speed, "s", 100, "speed in milliseconds")
	flag.Parse()

	if speed <= 0 {
		speed = 1
	}

	togglechan := make(chan bool)
	go listen(togglechan)
	clicker(togglechan)
}

func tapper(stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
			robotgo.Click("left")
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(speed)+1))
		}
	}
}

func listen(toggler chan bool) {
	var x bool
	var oldX bool
	stopChan := make(chan struct{}, 1)
	for {
		select {
		case x = <-toggler:
			if x {
				go tapper(stopChan)
			} else {
				if oldX {
					stopChan <- struct{}{}
				}
			}
			oldX = x
		}
	}
}

func clicker(toggler chan bool) {
	evChan := hook.Start()
	defer hook.End()

	var keyHolded bool

	for ev := range evChan {
		if ev.Rawcode == 65509 && ev.Kind == hook.KeyDown {
			keyHolded = true
		} else if ev.Rawcode == 65509 && ev.Kind == hook.KeyUp {
			toggler <- false
			keyHolded = false
		}

		if ev.Button == 3 && keyHolded {
			if ev.Kind == hook.MouseHold {
				toggler <- true
			} else if ev.Kind == hook.MouseDown {
				toggler <- false
			}
		}
	}
}
