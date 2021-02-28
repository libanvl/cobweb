package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dixonwille/wlog/v3"
	"github.com/dixonwille/wmenu/v5"
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

	cli := NewCli(*fbw, *fto)

	cliexepath, err := cli.CheckExePath()
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	ui.Info(fmt.Sprintf("Using cli: %s", cliexepath))
	fmt.Println()

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

	if status.Status != StatusUnlocked {
		ui.Error("BitWarden CLI must have an active unlocked session")
		os.Exit(3)
	}

	if *fsync {
		ui.Running("Syncing vault...")
		out, err := cli.Sync()
		if err != nil {
			if err == context.DeadlineExceeded {
				ui.Error("Operation timed out")
			} else {
				ui.Error(err.Error())
			}
			os.Exit(3)
		}
		ui.Success(out)
	} else {
		ui.Output("Skipping vault sync")
	}

	fmt.Println()

	menu := DoMenu(ui, "Which utility would you like to run?")
	menu.Option("Clean duplicates (exact)", func(ui wlog.UI, opt wmenu.Opt) error {
		return doDedupExact(ui, opt, cli)
	}, false, nil)
	menu.Option("Clean duplicates (fuzzy)", doDedupFuzzy, false, nil)
	menu.Option("Make URIs secure", doSecureUri, false, nil)
	err = menu.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func doSecureUri(wlog.UI, wmenu.Opt) error {
	return os.ErrInvalid
}

func doDedupFuzzy(wlog.UI, wmenu.Opt) error {
	return os.ErrInvalid
}

func printPreamble(ui wlog.UI) {
	ui.Output("Cobweb Utilities for BitWarden")
	ui.Output("==============================")
	fmt.Println()
}
