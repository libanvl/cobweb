package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dixonwille/wlog/v3"
	"github.com/libanvl/cobweb/cmd/cobweb/core"
	"github.com/libanvl/cobweb/cmd/cobweb/ext"
	"github.com/libanvl/cobweb/pkg/warden"
)

func main() {
	fbw := flag.String("bwpath", "bw", "The path to the bitwarden cli executable")
	fto := flag.Duration("timeout", 2*time.Second, "The timeout for cli operations")
	fsync := flag.Bool("sync", false, "Whether to sync the vault")
	fhelp := flag.Bool("help", false, "Show help for this command")
	flag.Parse()

	ui := DefaultUI()
	printPreamble(ui)

	if *fhelp {
		flag.Usage()
		os.Exit(0)
	}

	cli, err := warden.NewCli(*fbw, *fto)
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	ui.Info(fmt.Sprintf("Using cli: %s", cli.ExePath()))
	fmt.Println()

	checkStatus(ui, cli)

	if *fsync {
		ui.Running("Syncing vault...")
		out, err := cli.Sync()
		if err != nil {
			ui.Error(err.Error())
			os.Exit(3)
		}
		ui.Success(out)
	} else {
		ui.Output("Skipping vault sync")
	}

	fmt.Println()

	menubld := MenuBuilder{}

	opts := core.RunOpts{
		UI:          ui,
		MenuBuilder: menubld,
		Warden:      cli,
	}

	menu := menubld.DefaultMenu("Which utility would you like to run?")
	menu = menubld.AddRunRegistry(menu, &opts, ext.GlobalRunRegistry)
	menu = menubld.AddRunEntry(menu, &opts, helper{&opts})
	menu = menubld.AddExit(menu, 0, ui.Success, true)

	for {
		err = menu.Run()
		if err != nil && err != core.ErrQuitMenu {
			log.Fatal(err)
		}
		fmt.Println()
	}
}

func checkStatus(ui wlog.UI, cli *warden.Cli) {
	ui.Running("Getting status...")

	status, err := cli.Status()
	if err != nil {
		ui.Error(err.Error())
		os.Exit(2)
	}
	ui.Info(fmt.Sprintf("Server:\t%s", status.ServerURL))
	ui.Info(fmt.Sprintf("User:\t%s", status.UserEmail))
	ui.Info(fmt.Sprintf("Synced:\t%s", status.LastSync))
	ui.Info(fmt.Sprintf("Status:\t%s", status.Status))
	fmt.Println()

	if status.Status != warden.StatusUnlocked {
		ui.Error("BitWarden CLI must have an active unlocked session")
		ui.Info("$> export BW_SESSION=(bw --raw unlock)")
		os.Exit(3)
	}

	ui.Success("Done")
	fmt.Println()
}

func printPreamble(ui wlog.UI) {
	ui.Output("Cobweb Utilities for BitWarden")
	ui.Output("==============================")
	fmt.Println()
}
