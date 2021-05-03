package hosts

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/dixonwille/wmenu/v5"
	"github.com/libanvl/cobweb/cmd/cobweb/core"
	"github.com/libanvl/cobweb/pkg/warden"
)

type securehosts struct {
	opts  *core.RunOpts
	title string
}

func init() {
	var _ core.RunEntry = securehosts{}
	var _ core.HelpProvider = securehosts{}
	RunRegistry["Secure login hosts"] = func(ro *core.RunOpts, s string) (core.RunEntry, error) {
		return &securehosts{opts: ro, title: s}, nil
	}
}

type itemurl struct {
	item *warden.Item
	url  struct {
		index int
		url   *url.URL
	}
}

func (securehosts) Summary() string {
	return "Changes http URIs to https"
}

func (securehosts) Detail() string {
	return ""
}

func (s securehosts) Title() string {
	return s.title
}

func (s securehosts) Run() error {
	s.opts.UI.Running("Finding insecure hosts...")
	vault, err := s.opts.Warden.Vault()
	if err != nil {
		return err
	}

	filtered := []itemurl{}

	for _, item := range vault {
		if item.Type == warden.LoginItem {
			for x, uri := range item.Login.Uris {
				url, err := url.Parse(uri.URI)
				if err != nil {
					continue
				}
				if url.Scheme == "http" {
					r := itemurl{}
					r.item = item
					r.url.index = x
					r.url.url = url
					filtered = append(filtered, r)
				}
			}
		}
	}

	s.opts.UI.Success("Done")
	tofix := []itemurl{}
	lenfiltered := len(filtered)

	pagesize := 5
	numpage := (lenfiltered / pagesize)
	if lenfiltered%pagesize > 0 {
		numpage += 1
	}

	for x := 0; x < numpage; x++ {
		menu := s.opts.MenuBuilder.DefaultMenu(fmt.Sprintf("Select items to fix (page %d/%d): ", x+1, numpage))
		menu.AllowMultiple()
		menu.Action(func(opts []wmenu.Opt) error {
			for _, o := range opts {
				iu := o.Value.(itemurl)
				tofix = append(tofix, iu)
				s.opts.UI.Info(fmt.Sprintf("Added to fix list:\n%s", iu.item))
			}

			return nil
		})

		i := x * pagesize
		j := i + pagesize
		if lenfiltered < j {
			j = lenfiltered
		}

		for y := i; y < j; y++ {
			menu.Option(filtered[y].item.String(), &filtered[y], false, nil)
		}

		menu.Option("1-5", "PAGE", false, func(o wmenu.Opt) error {
			for y := i; y < j; y++ {
				iu := filtered[y]
				tofix = append(tofix, iu)
				s.opts.UI.Info(fmt.Sprintf("Added to fix list:\n%s", iu.item))
			}

			return nil
		})

		menu.Option("Process Fix List", "START", false, func(o wmenu.Opt) error {
			return errors.New("START")
		})

		if x < numpage-1 {
			menu.Option("Next Page", "NEXT", true, func(o wmenu.Opt) error {
				return nil
			})
		}

		menu = s.opts.MenuBuilder.AddQuitMenu(menu, x == numpage-1)

		if x == 0 {
			menu.Option("I'm feeling lucky", "LUCKY", false, func(o wmenu.Opt) error {
				for _, fiu := range filtered {
					localfiu := fiu
					tofix = append(tofix, localfiu)
					s.opts.UI.Info(fmt.Sprintf("Added to fix list:\n%s", localfiu.item))
				}
				s.opts.UI.Info(fmt.Sprintf("Added %d items to fix list", len(tofix)))

				return errors.New("START")
			})
		}

		err := menu.Run()
		if err != nil {
			if err.Error() == "START" {
				break
			} else {
				return err
			}
		}
	}

	err = updateFixList(s, tofix)
	if err != nil {
		return err
	}

	return nil
}

func updateFixList(s securehosts, tofix []itemurl) error {
	lentofix := len(tofix)
	if lentofix < 1 {
		return nil
	}

	menu := s.opts.MenuBuilder.DefaultMenu(fmt.Sprintf("Start updating %d items?", lentofix))
	menu.LoopOnInvalid()
	menu.IsYesNo(wmenu.DefY)
	menu.Action(func(o []wmenu.Opt) error {
		if o[0].ID == 0 {
			for x, iu := range tofix {
				url := iu.url.url
				url.Scheme = "https"

				item := iu.item
				item.Login.Uris[iu.url.index].URI = url.String()

				result, err := s.opts.Warden.EditItem(item)
				if err != nil {
					s.opts.UI.Error(err.Error())
				} else {
					s.opts.UI.Success(fmt.Sprintf("Updated item %d/%d: %s", x+1, lentofix, result))
				}
			}

			s.opts.UI.Success("Done")
			return nil
		}

		s.opts.UI.Warn("Skipping updates")
		return nil
	})

	if err := menu.Run(); err != nil {
		return err
	}

	return nil
}
