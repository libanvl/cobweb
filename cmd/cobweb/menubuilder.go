package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dixonwille/wlog/v3"
	"github.com/dixonwille/wmenu/v5"
	"github.com/libanvl/cobweb/cmd/cobweb/ext"
)

func init() {
	var _ ext.MenuBuilder = MenuBuilder{}
}

type MenuBuilder struct {
}

func (MenuBuilder) DefaultMenu(question string) *wmenu.Menu {
	menu := wmenu.NewMenu(question)
	menu.AddColor(wlog.BrightBlue, wlog.Green, wlog.Cyan, wlog.Magenta)
	return menu
}

func (fac MenuBuilder) AddRunEntries(menu *wmenu.Menu, opts *ext.RunOpts, entries ext.RunRegistry) *wmenu.Menu {
	runfunc := func(o wmenu.Opt) error {
		entryfac := o.Value.(ext.RunEntryFactory)

		entry, err := entryfac(opts, o.Text)
		if err != nil {
			return err
		}

		fmt.Println()
		opts.UI.Output(o.Text)
		opts.UI.Output(strings.Repeat("=", len(o.Text)))
		fmt.Println()

		if err = entry.Run(); err != nil {
			return err
		}

		return nil
	}

	for title, entryfac := range entries {
		menu.Option(title.String(), entryfac, false, runfunc)
	}

	return menu
}

func (MenuBuilder) AddExit(menu *wmenu.Menu, code int, printer func(string), isDefault bool) *wmenu.Menu {
	menu.Option("Exit", "", isDefault, func(o wmenu.Opt) error {
		printer("exiting")
		os.Exit(code)
		return nil
	})

	return menu
}
