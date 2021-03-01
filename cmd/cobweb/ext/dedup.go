package ext

import (
	"fmt"

	"github.com/dixonwille/wmenu/v5"
	"github.com/libanvl/cobweb/pkg/warden"
)

func init() {
	var _ RunEntry = dedup{}
	GlobalRunRegistry["Clean Duplicate Hostnames (Exact)"] = func(ro *RunOpts, t string) (RunEntry, error) {
		return dedup{opts: ro, title: t}, nil
	}
}

type dedup struct {
	opts  *RunOpts
	title string
}

func (d dedup) Run() error {
	ui := d.opts.UI
	cli := d.opts.Warden

	ui.Running("Finding duplicates")
	vault, err := cli.Vault()
	if err != nil {
		ui.Error(string(err.Error()))
		return err
	}

	hnitems := vault.HostnameItemMap(func(i *warden.Item, s string, e error) {
		ui.Warn(fmt.Sprintf("ID: %s, URL: %s", i.ID, s))
		ui.Warn(e.Error())
	})

	for hname, items := range hnitems {
		if len(items) < 2 {
			delete(hnitems, hname)
		}
	}

	ui.Success(fmt.Sprintf("Found duplicates: %d", len(hnitems)))
	fmt.Println()

	return d.processMap(hnitems)
}

func (d *dedup) processMap(hnitems warden.HostnameItemMap) error {
	ui := d.opts.UI
	cli := d.opts.Warden
	menubld := d.opts.MenuBuilder

	for hname, items := range hnitems {
		ui.Info(hname)

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
			return err
		}
	}

	return nil
}
