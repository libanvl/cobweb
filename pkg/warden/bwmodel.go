package warden

import (
	"fmt"
	"time"
)

type Status struct {
	ServerURL string    `json:"serverUrl"`
	LastSync  time.Time `json:"lastSync"`
	UserEmail string    `json:"userEmail"`
	UserID    string    `json:"userId"`
	Status    string    `json:"status"`
}

const (
	StatusUnauthenticated string = "unauthenticated"
	StatusLocked          string = "locked"
	StatusUnlocked        string = "unlocked"
)

type Item struct {
	Object         string        `json:"object"`
	ID             string        `json:"id"`
	OrganizationID interface{}   `json:"organizationId"`
	FolderID       interface{}   `json:"folderId"`
	Type           ItemType      `json:"type"`
	Name           string        `json:"name"`
	Notes          interface{}   `json:"notes"`
	Favorite       bool          `json:"favorite"`
	Login          Login         `json:"login"`
	CollectionIds  []interface{} `json:"collectionIds"`
	RevisionDate   time.Time     `json:"revisionDate"`
}

//go:generate golang.org/x/tools/cmd/stringer -type=ItemType
type ItemType int

const (
	LoginItem  ItemType = 1
	SecureNote ItemType = 2
	Card       ItemType = 3
	Identity   ItemType = 4
)

func (item Item) String() string {
	return fmt.Sprintf("%s %s %s", item.Name, item.RevisionDate.Format(time.UnixDate), item.Login)
}

type Uri struct {
	Match MatchType `json:"match"`
	URI   string    `json:"uri"`
}

type Uris []Uri

func (uris Uris) String() string {
	return uris[0].URI
}

//go:generate golang.org/x/tools/cmd/stringer -type=MatchType
type MatchType int

const (
	Domain     MatchType = 0
	Host       MatchType = 1
	StartsWith MatchType = 2
	Exact      MatchType = 3
	Regular    MatchType = 4
	Never      MatchType = 5
)

type Login struct {
	Uris                 Uris        `json:"uris"`
	Username             string      `json:"username"`
	Password             string      `json:"password"`
	Totp                 interface{} `json:"totp"`
	PasswordRevisionDate time.Time   `json:"passwordRevisionDate"`
}

func (login Login) String() string {
	uris := ""
	for _, u := range login.Uris {
		uris += fmt.Sprintf(" %s\n", u.URI)
	}

	return fmt.Sprintf("%s %s\n%s", login.Username, login.Password, uris)
}
