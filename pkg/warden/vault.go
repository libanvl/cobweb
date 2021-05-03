package warden

import (
	"net/url"
)

type Vault []*Item

type HostnameItemMap map[string][]*Item
type PasswordItemMap map[string][]*Item

func (v Vault) HostnameItemMap(onerr func(*Item, string, error)) HostnameItemMap {
	result := make(HostnameItemMap, 0)
	for _, item := range v {
		if item.Type != LoginItem {
			continue
		}

		for _, bwuri := range item.Login.Uris {
			url, err := url.Parse(bwuri.URI)
			if err != nil {
				if onerr != nil {
					onerr(item, bwuri.URI, err)
				}

				continue
			}

			hname := url.Hostname()
			if _, ok := result[hname]; !ok {
				result[hname] = []*Item{item}
			} else {
				result[hname] = append(result[hname], item)
			}
		}
	}

	return result
}

func (v Vault) PasswordItemMap() PasswordItemMap {
	result := make(PasswordItemMap, 0)
	for _, item := range v {
		if item.Type != LoginItem {
			continue
		}

		pword := item.Login.Password
		if _, ok := result[pword]; !ok {
			result[pword] = []*Item{item}
		} else {
			result[pword] = append(result[pword], item)
		}
	}

	return result
}

func (v Vault) Filter(filter func(*Item) bool) Vault {
	result := make(Vault, 0)
	for _, item := range v {
		if !filter(item) {
			continue
		}

		result = append(result, item)
	}

	return result
}
