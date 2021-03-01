package ext

import (
	"github.com/dixonwille/wlog/v3"
	"github.com/dixonwille/wmenu/v5"
	"github.com/libanvl/cobweb/pkg/warden"
)

var GlobalRunRegistry RunRegistry = make(RunRegistry, 0)

type RunRegistry map[EntryTitle]RunEntryFactory

type RunEntryFactory func(*RunOpts, string) (RunEntry, error)

type EntryTitle string

type MenuBuilder interface {
	DefaultMenu(string) *wmenu.Menu
	AddRunEntries(*wmenu.Menu, *RunOpts, RunRegistry) *wmenu.Menu
	AddExit(*wmenu.Menu, int, func(string), bool) *wmenu.Menu
}

type RunOpts struct {
	UI          wlog.UI
	MenuBuilder MenuBuilder
	Warden      warden.Warden
}

type RunEntry interface {
	Run() error
}

func (et EntryTitle) String() string {
	return string(et)
}
