package passdup

import (
	"fmt"

	"github.com/libanvl/cobweb/cmd/cobweb/core"
)

type pod struct {
	opts  *core.RunOpts
	title string
}

func init() {
	var _ core.RunEntry = pod{}
	RunRegistry["Print duplicate passwords (Exact)"] = func(ro *core.RunOpts, s string) (core.RunEntry, error) {
		return &pod{opts: ro, title: s}, nil
	}
}

func (p pod) Title() string {
	return p.title
}

func (p pod) Run() error {
	vault, err := p.opts.Warden.Vault()
	if err != nil {
		return err
	}

	pitems := vault.PasswordItemMap()
	for pword, items := range pitems {
		if len(items) < 2 {
			continue
		}

		p.opts.UI.Info(pword)
		for x, item := range items {
			p.opts.UI.Output(fmt.Sprintf("%4d %s", x, item))
		}
	}

	return nil
}
