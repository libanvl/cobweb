package ext

import (
	"fmt"

	"github.com/libanvl/cobweb/cmd/cobweb/core"
	"github.com/libanvl/cobweb/cmd/cobweb/ext/passdup"
)

type extpassdup struct {
	opts  *core.RunOpts
	title string
}

func init() {
	var _ core.RunEntry = extpassdup{}
	var _ core.HelpProvider = extpassdup{}
	GlobalRunRegistry["Duplicate Passwords"] = func(ro *core.RunOpts, s string) (core.RunEntry, error) {
		return &extpassdup{ro, s}, nil
	}
}

func (p extpassdup) Summary() string {
	return "Utilities for handling duplicate passwords"
}

func (p extpassdup) Detail() string {
	return ""
}

func (p extpassdup) Title() string {
	return p.title
}

func (p extpassdup) Run() error {
	menu := p.opts.MenuBuilder.DefaultMenu("Which password utility would you like to run?")
	menu = p.opts.MenuBuilder.AddRunRegistry(menu, p.opts, passdup.RunRegistry)
	menu = p.opts.MenuBuilder.AddQuitMenu(menu, true)

	for {
		err := menu.Run()
		if err != nil {
			return err
		}
		fmt.Println()
	}
}
