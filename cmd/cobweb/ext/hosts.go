package ext

import (
	"fmt"

	"github.com/libanvl/cobweb/cmd/cobweb/core"
	"github.com/libanvl/cobweb/cmd/cobweb/ext/hosts"
)

type exthosts struct {
	opts  *core.RunOpts
	title string
}

func init() {
	var _ core.RunEntry = exthosts{}
	var _ core.HelpProvider = exthosts{}
	GlobalRunRegistry["Hosts Utilities"] = func(ro *core.RunOpts, s string) (core.RunEntry, error) {
		return &exthosts{opts: ro, title: s}, nil
	}
}

func (h exthosts) Summary() string {
	return "Utilities for handling login URIs"
}

func (h exthosts) Detail() string {
	return ""
}

func (h exthosts) Title() string {
	return h.title
}

func (h exthosts) Run() error {
	menu := h.opts.MenuBuilder.DefaultMenu("Which utilility would you like to run?")
	menu = h.opts.MenuBuilder.AddRunRegistry(menu, h.opts, hosts.RunRegistry)
	menu = h.opts.MenuBuilder.AddQuitMenu(menu, true)

	for {
		err := menu.Run()
		if err != nil {
			return err
		}
		fmt.Println()
	}
}
