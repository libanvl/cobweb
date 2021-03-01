package main

import (
	"os"

	"github.com/dixonwille/wlog/v3"
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
