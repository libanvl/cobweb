package main

import (
	"fmt"
	"os"

	"github.com/dixonwille/wlog/v3"
	"github.com/dixonwille/wmenu/v5"
)

func doDedupExact(ui wlog.UI, opt wmenu.Opt, cli *Cli) error {
	ui.Output("Clean Duplicates (Exact)")
	ui.Output("========================")
	fmt.Println()

	ui.Running("Finding duplicates")
	vault, err := cli.Vault()
	if err != nil {
		ui.Error(string(err.Error()))
		return err
	}

	hnitems := vault.HostnameItemMap(func(i *Item, s string, e error) {
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

	return processMap(ui, hnitems, cli)
}

func processMap(ui wlog.UI, hnitems HostnameItemMap, cli *Cli) error {
	for hname, items := range hnitems {
		ui.Info(hname)

		menu := DefaultMenu("Delete which items?")
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
				item, ok := opt.Value.(*Item)
				if !ok {
					panic("Unexpected type")
				}

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
		menu.Option("Exit", "EXIT", false, func(o wmenu.Opt) error {
			ui.Info("Exiting")
			os.Exit(0)
			return nil
		})

		err := menu.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
