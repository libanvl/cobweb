package core

import (
	"errors"

	"github.com/dixonwille/wlog/v3"
	"github.com/dixonwille/wmenu/v5"
	"github.com/libanvl/cobweb/pkg/warden"
)

var ErrQuitMenu error = errors.New("Quit current menu")

type RunRegistry map[EntryTitle]RunEntryFactory

type RunEntryFactory func(*RunOpts, string) (RunEntry, error)

type EntryTitle string

type MenuBuilder interface {
	DefaultMenu(string) *wmenu.Menu
	AddRunEntry(*wmenu.Menu, *RunOpts, RunEntry) *wmenu.Menu
	AddRunRegistry(*wmenu.Menu, *RunOpts, RunRegistry) *wmenu.Menu
	AddQuitMenu(menu *wmenu.Menu, isDefault bool) *wmenu.Menu
}

type RunOpts struct {
	UI          wlog.UI
	MenuBuilder MenuBuilder
	Warden      *warden.Cli
}

type RunEntry interface {
	Title() string
	Run() error
}

type HelpProvider interface {
	Summary() string
	Detail() string
}

func (et EntryTitle) String() string {
	return string(et)
}
