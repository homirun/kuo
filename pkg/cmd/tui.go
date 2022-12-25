package cmd

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

type firstColumn struct {
	tuiMode
	context       string
	kubectlStdOut string
}
type secondColumn struct {
	tuiMode
	context       string
	kubectlStdOut string
}

type tuiMode struct {
	mode string
}

func ShowTui(kubectlStdOutputMaps map[string]string, tuiMode string) error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	var contexts []string
	var outputs []string
	for context, output := range kubectlStdOutputMaps {
		contexts = append(contexts, context)
		outputs = append(outputs, output)
	}

	fc := new(firstColumn)
	sc := new(secondColumn)

	fc.tuiMode.mode = tuiMode
	fc.context = contexts[0]
	fc.kubectlStdOut = outputs[0]
	sc.tuiMode.mode = tuiMode
	sc.context = contexts[1]
	sc.kubectlStdOut = outputs[1]

	g.SetManager(fc, sc)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func (fc *firstColumn) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if fc.tuiMode.mode == "h" {
		if v, err := g.SetView("firstColumn", 0, 0, maxX-1, maxY/2-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			fmt.Fprintln(v, fc.kubectlStdOut)
		}
	} else {
		if v, err := g.SetView("firstColumn", 0, 0, maxX/2-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			fmt.Fprintln(v, fc.kubectlStdOut)
		}
	}

	return nil
}

func (sc *secondColumn) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if sc.tuiMode.mode == "h" {
		if v, err := g.SetView("SecondColumn", 0, maxY/2, maxX-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			fmt.Fprintln(v, sc.kubectlStdOut)
		}
	} else {
		if v, err := g.SetView("SecondColumn", maxX/2, 0, maxX-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			fmt.Fprintln(v, sc.kubectlStdOut)
		}
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
