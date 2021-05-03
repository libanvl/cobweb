package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dixonwille/wlog/v3"
	"github.com/dixonwille/wmenu/v5"
	"github.com/libanvl/cobweb/cmd/cobweb/core"
)

func init() {
	var _ core.MenuBuilder = MenuBuilder{}
}

type MenuBuilder struct {
}

func (MenuBuilder) DefaultMenu(question string) *wmenu.Menu {
	menu := wmenu.NewMenu(question)
	menu.LoopOnInvalid()
	menu.AddColor(wlog.BrightBlue, wlog.Green, wlog.Cyan, wlog.Magenta)
	return menu
}

func (MenuBuilder) AddRunEntry(menu *wmenu.Menu, opts *core.RunOpts, entry core.RunEntry) *wmenu.Menu {
	runfunc := func(o wmenu.Opt) error {
		entry := o.Value.(core.RunEntry)

		fmt.Println()
		opts.UI.Output(entry.Title())
		opts.UI.Output(strings.Repeat("=", len(entry.Title())))
		fmt.Println()

		if err := entry.Run(); err != nil {
			return err
		}

		return nil
	}

	title := entry.Title()

	if helper, ok := entry.(core.HelpProvider); ok {
		if summary := helper.Summary(); summary != "" {
			title = fmt.Sprintf("%s\n%s", title, summary)
		}
	}

	title += "\n"

	menu.Option(title, entry, false, runfunc)

	return menu
}

func (fac MenuBuilder) AddRunRegistry(menu *wmenu.Menu, opts *core.RunOpts, entries core.RunRegistry) *wmenu.Menu {
	for et, entryfac := range entries {
		title := et.String()
		entry, err := entryfac(opts, title)
		if err != nil {
			continue
		}

		menu = fac.AddRunEntry(menu, opts, entry)
	}

	return menu
}

func (MenuBuilder) AddQuitMenu(menu *wmenu.Menu, isDefault bool) *wmenu.Menu {
	menu.Option("Quit", "", isDefault, func(o wmenu.Opt) error {
		return core.ErrQuitMenu
	})

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
