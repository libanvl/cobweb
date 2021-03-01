package ext

import (
	"fmt"

	"github.com/dixonwille/wmenu/v5"
	"github.com/libanvl/cobweb/pkg/warden"
)

func init() {
	var _ RunEntry = &passdedup{}
	GlobalRunRegistry["Clean Duplicate Passwords (Exact)"] = func(ro *RunOpts, t string) (RunEntry, error) {
		return passdedup{opts: ro, title: t}, nil
	}
}

type passdedup struct {
	opts  *RunOpts
	title string
}

func (p passdedup) Run() error {
	ui := p.opts.UI
	cli := p.opts.Warden

	ui.Running("Finding duplicates")
	vault, err := cli.Vault()
	if err != nil {
		ui.Error(string(err.Error()))
		return err
	}

	pworditems := vault.PasswordItemMap()

	for hname, items := range pworditems {
		if len(items) < 2 {
			delete(pworditems, hname)
		}
	}

	ui.Success(fmt.Sprintf("Found duplicates: %d", len(pworditems)))
	fmt.Println()

	return p.processPasswordMap(pworditems)
}

func (p *passdedup) processPasswordMap(pworditems warden.PasswordItemMap) error {
	ui := p.opts.UI
	cli := p.opts.Warden
	menubld := p.opts.MenuBuilder

	for pword, items := range pworditems {
		ui.Info(pword)

		menu := menubld.DefaultMenu("Delete which items?")
		menu.AllowMultiple()
		menu.LoopOnInvalid()
		menu.Action(func(opts []wmenu.Opt) error {
			if opts[0].Text == "Skip" {
				return nil
			}

			if opts[0].Text == "" {
				return wmenu.ErrNoResponse
			}

			ui.Running(fmt.Sprintf("Deleting items"))
			for _, opt := range opts {
				item := opt.Value.(*warden.Item)
				if err := cli.DeleteItem(item); err != nil {
					ui.Error(err.Error())
				}

				ui.Success(fmt.Sprintf("Deleted item: %s", item))
			}

			return nil
		})

		for _, item := range items {
			menu.Option(fmt.Sprintf("%s", item), item, false, nil)
		}

		menu.Option("Skip", "SKIP", true, nil)
		menu = menubld.AddExit(menu, 0, ui.Success, false)
		err := menu.Run()
		if err != nil {
			ui.Error(err.Error())
		}
	}

	return nil
}
