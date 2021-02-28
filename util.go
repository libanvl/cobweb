package main

import (
	"errors"
	"os"

	"github.com/dixonwille/wlog/v3"
	"github.com/dixonwille/wmenu/v5"
)

func DefaultUI() wlog.UI {
	var ui wlog.UI = wlog.New(os.Stdin, os.Stdout, os.Stderr)
	ui = wlog.AddPrefix("?", wlog.Cross, "", "", "", "=>", wlog.Check, "!", ui)
	ui = wlog.AddColor(
		wlog.Cyan,
		wlog.Red,
		wlog.BrightMagenta,
		wlog.BrightBlue,
		wlog.Green,
		wlog.None,
		wlog.None,
		wlog.BrightGreen,
		wlog.Yellow,
		ui)

	return ui
}

func DefaultMenu(question string) *wmenu.Menu {
	menu := wmenu.NewMenu(question)
	menu.LoopOnInvalid()
	menu.AddColor(wlog.BrightBlue, wlog.Green, wlog.Cyan, wlog.Magenta)
	return menu
}

func DoMenu(ui wlog.UI, question string) *wmenu.Menu {
	menu := DefaultMenu(question)
	menu.Action(func(opts []wmenu.Opt) error {
		for _, opt := range opts {
			do, ok := opt.Value.(func(wlog.UI, wmenu.Opt) error)
			if !ok {
				return errors.New("Internal error")
			}
			if err := do(ui, opt); err != nil {
				return err
			}
		}

		return nil
	})

	menu.Option("Exit", func(ui wlog.UI, _ wmenu.Opt) error {
		ui.Success("exiting")
		os.Exit(0)
		return nil
	}, true, nil)

	return menu
}
