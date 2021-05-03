package hosts

import (
	"fmt"

	"github.com/dixonwille/wmenu/v5"
	"github.com/libanvl/cobweb/cmd/cobweb/core"
	"github.com/libanvl/cobweb/pkg/warden"
)

func init() {
	var _ core.RunEntry = dedup{}
	var _ core.HelpProvider = dedup{}
	RunRegistry["Clean Duplicate Hostnames (Exact)"] = func(ro *core.RunOpts, t string) (core.RunEntry, error) {
		return dedup{opts: ro, title: t}, nil
	}
}

type dedup struct {
	opts  *core.RunOpts
	title string
}

func (d dedup) Title() string {
	return d.title
}

func (dedup) Summary() string {
	return "Groups login items by hostname. Select items to delete from each group."
}

func (dedup) Detail() string {
	return "Matches are exact string matches of the hostname component of the login URIs"
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
		menu = menubld.AddQuitMenu(menu, false)
		err := menu.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
