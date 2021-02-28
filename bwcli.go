package main

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"time"
)

type Cli struct {
	bwexe   string
	timeout time.Duration
}

type operation func(context.Context) error

func NewCli(bwexe string, timeout time.Duration) *Cli {
	return &Cli{
		bwexe:   bwexe,
		timeout: timeout,
	}
}

func (cli *Cli) CheckExePath() (string, error) {
	bwexe, err := exec.LookPath(cli.bwexe)
	if err != nil {
		return "", err
	}

	cli.bwexe = bwexe
	return cli.bwexe, nil
}

func (cli Cli) Sync() (string, error) {
	var out []byte
	err := cli.withTimeout(
		context.Background(),
		func(ctx context.Context) error {
			cmd := cli.command(ctx, "sync")
			var _err error
			out, _err = cmd.CombinedOutput()
			if _err != nil {
				return errors.New(string(out))
			}

			return nil
		})

	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (cli Cli) Status() (*Status, error) {
	var out []byte
	var status Status
	err := cli.withTimeout(
		context.Background(),
		func(ctx context.Context) error {
			cmd := cli.command(ctx, "status", "--raw")
			var _err error
			out, _err = cmd.Output()
			if _err != nil {
				return _err
			}

			return nil
		})

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(out, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func (cli Cli) Vault() (Vault, error) {
	var out []byte
	var vault Vault
	err := cli.withTimeout(
		context.Background(),
		func(ctx context.Context) error {
			cmd := cli.command(ctx, "list", "items", "--raw")
			var _err error
			out, _err = cmd.Output()
			if _err != nil {
				return _err
			}
			return nil
		})

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(out, &vault); err != nil {
		return nil, err
	}

	return vault, nil
}

func (cli Cli) DeleteItem(item *Item) error {
	return cli.withTimeout(
		context.Background(),
		func(ctx context.Context) error {
			cmd := cli.command(ctx, "delete", "item", item.ID)
			return cmd.Run()
		})
}

func (cli Cli) command(ctx context.Context, arg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, cli.bwexe, arg...)
}

func (cli Cli) withTimeout(parent context.Context, op operation) error {
	if op == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(parent, cli.timeout)
	defer cancel()
	cmderr := op(ctx)
	if cxterr := ctx.Err(); cxterr != nil {
		return cxterr
	}

	if cmderr != nil {
		if exerr, ok := cmderr.(*exec.ExitError); ok {
			return errors.New(string(exerr.Stderr))
		} else {
			return cmderr
		}
	}

	return nil
}
