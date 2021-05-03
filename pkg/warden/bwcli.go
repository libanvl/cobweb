package warden

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

type Cli struct {
	bwexe   string
	timeout time.Duration
}

type CliError struct {
	cli    *Cli
	output string
}

func (e *CliError) Error() string {
	return fmt.Sprintf("cli error: [%s]\n%s", e.cli.bwexe, e.output)
}

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
	var err error

	err = c.withTimeout(
		context.Background(),
		func(ctx context.Context) error {
			out, err = c.command(ctx, "sync").CombinedOutput()
			if err != nil {
				return c.newError(string(out))
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
	var err error
	var status Status

	err = c.withTimeout(
		context.Background(),
		func(ctx context.Context) error {
			out, err = c.command(ctx, "status", "--raw").CombinedOutput()
			if err != nil {
				return c.newError(string(out))
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
	var err error
	var vault Vault

	err = c.withTimeout(
		context.Background(),
		func(ctx context.Context) error {
			out, err = c.command(ctx, "list", "items", "--raw").CombinedOutput()
			if err != nil {
				return c.newError(string(out))
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
			out, err := c.command(ctx, "delete", "item", item.ID).CombinedOutput()
			if err != nil {
				return c.newError(string(out))
			}

			return nil
		})
}

func (c Cli) EditItem(item *Item) (*Item, error) {
	var out []byte
	var err error
	var result Item

	err = c.withTimeout(
		context.Background(),
		func(ctx context.Context) error {
			updated, err := EncodeItem(item)
			if err != nil {
				return err
			}

			out, err = c.command(ctx, "edit", "item", item.ID, updated).CombinedOutput()
			if err != nil {
				return c.newError(string(out))
			}
			return nil
		})

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(out, &result); err != nil {
		return nil, err
	}

	return &result, err
}

func EncodeItem(item *Item) (string, error) {
	json, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(json), nil
}

func (c Cli) newError(output string) *CliError {
	return &CliError{cli: &c, output: output}
}

func (c Cli) command(ctx context.Context, arg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, c.bwexe, arg...)
}

func (c Cli) withTimeout(parent context.Context, op func(context.Context) error) error {
	if op == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(parent, c.timeout)
	defer cancel()

	cmderr := op(ctx)
	if ctxerr := ctx.Err(); ctxerr != nil {
		return ctxerr
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
