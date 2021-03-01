package warden

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"time"
)

type Warden interface {
	ExePath() string
	Sync() (string, error)
	Status() (*Status, error)
	Vault() (Vault, error)
	DeleteItem(*Item) error
}

type cli struct {
	bwexe   string
	timeout time.Duration
}

type operation func(context.Context) error

func init() {
	var _ Warden = cli{}
}

func NewWarden(bwexe string, timeout time.Duration) (Warden, error) {
	bwexe, err := exec.LookPath(bwexe)
	if err != nil {
		return nil, err
	}

	return &cli{bwexe: bwexe, timeout: timeout}, nil
}

func (c cli) ExePath() string {
	return c.bwexe
}

func (c cli) Sync() (string, error) {
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

func (c cli) Status() (*Status, error) {
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

func (c cli) Vault() (Vault, error) {
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

func (c cli) DeleteItem(item *Item) error {
	return c.withTimeout(
		context.Background(),
		func(ctx context.Context) error {
			cmd := c.command(ctx, "delete", "item", item.ID)
			return cmd.Run()
		})
}

func (c cli) command(ctx context.Context, arg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, c.bwexe, arg...)
}

func (c cli) withTimeout(parent context.Context, op operation) error {
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
