package warden

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

func NewCli(bwexe string, timeout time.Duration) (*Cli, error) {
	bwexe, err := exec.LookPath(bwexe)
	if err != nil {
		return nil, err
	}

	return &Cli{bwexe: bwexe, timeout: timeout}, nil
}

func (c Cli) ExePath() string {
	return c.bwexe
}

func (c Cli) Sync() (string, error) {
	var out []byte
	err := c.withTimeout(
		context.Background(),
		func(ctx context.Context) error {
			cmd := c.command(ctx, "sync")
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

func (c Cli) Status() (*Status, error) {
	var out []byte
	var status Status
	err := c.withTimeout(
		context.Background(),
		func(ctx context.Context) error {
			cmd := c.command(ctx, "status", "--raw")
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

func (c Cli) Vault() (Vault, error) {
	var out []byte
	var vault Vault
	err := c.withTimeout(
		context.Background(),
		func(ctx context.Context) error {
			cmd := c.command(ctx, "list", "items", "--raw")
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

func (c Cli) DeleteItem(item *Item) error {
	return c.withTimeout(
		context.Background(),
		func(ctx context.Context) error {
			cmd := c.command(ctx, "delete", "item", item.ID)
			return cmd.Run()
		})
}

func (c Cli) command(ctx context.Context, arg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, c.bwexe, arg...)
}

func (c Cli) withTimeout(parent context.Context, op operation) error {
	if op == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(parent, c.timeout)
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
