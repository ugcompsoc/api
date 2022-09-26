package services

import (
	"crypto/tls"

	ldap "github.com/go-ldap/ldap/v3"
	"github.com/nuigcompsoc/api/internal/config"
	log "github.com/sirupsen/logrus"
)

type LdapService struct {
	Conn              *ldap.Conn
	Bind              string
	Password          string
	UserSearchBase    string
	SocietySearchBase string
	GroupSearchBase   string
	SearchBase        string
}

func NewLdap(config *config.Config) *LdapService {
	l, err := ldap.DialURL(config.LDAP.URL)
	if err != nil {
		log.WithField("error", err).Fatal("Failed to connect to LDAP Service")
	}

	// Reconnect with TLS
	err = l.StartTLS(&tls.Config{})
	if err != nil {
		log.WithField("error", err).Fatal("Failed to upgrade LDAP connection to TLS")
	}

	return &LdapService{
		Conn:              l,
		Bind:              config.LDAP.Bind,
		Password:          config.LDAP.Password,
		UserSearchBase:    config.LDAP.UserSearchBase,
		SocietySearchBase: config.LDAP.SocietySearchBase,
		GroupSearchBase:   config.LDAP.GroupSearchBase,
		SearchBase:        config.LDAP.SearchBase,
	}
}
