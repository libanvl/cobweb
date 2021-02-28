package main

import (
	"net/url"
)

type Vault []*Item

type HostnameItemMap map[string][]*Item

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
