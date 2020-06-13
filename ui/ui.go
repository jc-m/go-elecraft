package ui

import (
	"fmt"
	"math"

	"github.com/awesome-gocui/gocui"
	"github.com/w6ipa/go-elecraft/utils"
)

func CWPracticeLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	maxTopY := 0 + math.Round(float64(maxY/2))
	maxBottomY := maxTopY + math.Round(float64(maxY/2))

	if v, err := g.SetView("top", 1, 1, maxX-1, int(maxTopY), 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Wrap = true
		if _, err := g.SetCurrentView("top"); err != nil {
			return err
		}
		g.CurrentView().Title = "Active"
	}
	if v, err := g.SetView("bottom", 1, int(maxTopY+1), maxX-1, int(maxBottomY-1), 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Wrap = true
		v.Autoscroll = true
	}
	return nil
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func ScrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}

func CWUpdate(g *gocui.Gui, c chan []byte, done chan struct{}) {
Loop:
	for {
		select {
		case <-done:
			return
		case data, ok := <-c:
			if !ok {
				break Loop
			}
			g.Update(func(g *gocui.Gui) error {
				bottom, err := g.View("bottom")
				if err != nil {
					return err
				}
				fmt.Fprintf(bottom, "%s", data)
				top, err := g.View("top")
				if err != nil {
					return err
				}
				x, y := top.Cursor()

				line, err := top.Line(y)
				if err != nil {
					return err
				}
				dx := utils.CheckAndAdvance([]byte(line), x, data)
				top.MoveCursor(dx, 0, false)
				return nil
			})
		}
	}
	return
}
