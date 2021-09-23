package ldap

import (
	"github.com/go-ldap/ldap/v3"
	"github.com/nuigcompsoc/api/internal/config"
)

func bind(c *config.Config) (*ldap.Conn) {
	l, err := ldap.DialURL(c.LDAP.URL)
	if err != nil {
		return nil
	}

	err = l.Bind(c.LDAP.BindUser, c.LDAP.BindSecret)
	if err != nil {
		return nil
	}
	return l
}