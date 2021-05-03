package main

import (
	"fmt"
	"strings"

	"github.com/dixonwille/wmenu/v5"
	"github.com/libanvl/cobweb/cmd/cobweb/core"
	"github.com/libanvl/cobweb/cmd/cobweb/ext"
)

type helper struct {
	opts *core.RunOpts
}

func init() {
	var _ core.RunEntry = helper{}
}

func (helper) Title() string {
	return "Extension Help"
}

func (h helper) Run() error {
	menu := h.opts.MenuBuilder.DefaultMenu("Show help for which extension?")
	for et, entryfac := range ext.GlobalRunRegistry {
		entry, err := entryfac(h.opts, et.String())
		if err != nil {
			return err
		}

		if helper, ok := entry.(core.HelpProvider); ok {
			menu.Option(et.String(), helper, false, func(o wmenu.Opt) error {
				e := o.Value.(core.HelpProvider)
				title := o.Text
				lentitle := len(title)

				fmt.Println()
				h.opts.UI.Info(o.Text)
				h.opts.UI.Info(strings.Repeat("-", lentitle))
				h.opts.UI.Info(e.Summary())
				fmt.Println()
				h.opts.UI.Info(e.Detail())
				h.opts.UI.Info(strings.Repeat("=", lentitle))
				fmt.Println()
				return nil
			})
		}
	}

	menu = h.opts.MenuBuilder.AddQuitMenu(menu, true)
	for {
		err := menu.Run()
		if err != nil {
			return err
		}
	}
}
