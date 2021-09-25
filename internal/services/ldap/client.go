package ldap

import (
	"github.com/go-ldap/ldap/v3"
	"github.com/nuigcompsoc/api/internal/config"
)

type User struct {
	FullName 	string
	FirstName	string
	Surname		string
	MemberID	string
	Mail		[]string
	UID			string
	ObjectClass	[]string
	Shell		string
	Home		string
	UIDNumber	int
}

type Client struct {
	URL			string
	DN			string
	BindUser	string
	BindSecret	string
	Attributes	[]string
}

func NewClient(c *config.Config) *Client {
	return &Client{
		URL: c.LDAP.URL,
		DN: c.LDAP.DN,
		BindUser: c.LDAP.BindUser,
		BindSecret: c.LDAP.BindSecret,
		Attributes: c.LDAP.Attributes,
	}
}

func (c *Client) bind() (*ldap.Conn) {
	l, err := ldap.DialURL(c.URL)
	if err != nil {
		return nil
	}

	err = l.Bind(c.BindUser, c.BindSecret)
	if err != nil {
		return nil
	}
	return l
}