package ldap

import (
	log "github.com/sirupsen/logrus"
	"github.com/go-ldap/ldap/v3"
	"github.com/nuigcompsoc/api/internal/config"
	"strings"
)

/*
 * LDAP Utils
 */

 func entriesToMap(entries []*ldap.Entry, attributes []string) (map[string]map[string][]string) {
	arr := make(map[string]map[string][]string, len(entries))
	for _, entry := range entries {
		arr[entry.GetAttributeValue("uid")] = make(map[string][]string, len(attributes))
		for _, attribute := range attributes {
			values := entry.GetAttributeValues(attribute)
			arr[entry.GetAttributeValue("uid")][attribute] = make([]string, len(values))
			arr[entry.GetAttributeValue("uid")][attribute] = values
		}
	}
	return arr
}
 
/*
 * Account Search Utils
 */

func search(c *config.Config, searchBase string, filter string, attributes []string) ([]*ldap.Entry, bool) {
	l := bind(c)
	searchRequest := ldap.NewSearchRequest(
		searchBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.WithFields(log.Fields{
			"message": "error while searching",
			"error": err.Error(),
		}).Error("ldap")
		return nil, false
	}
	
	return sr.Entries, true
}

func GetUsersOrSocieties(c *config.Config, group string, attributes []string) (map[string]map[string][]string, bool) {
	return GetUserOrSociety(c, "*", group, attributes)
}

func GetUserOrSociety(c *config.Config, uid string, group string, attributes []string) (map[string]map[string][]string, bool) {
	entries, ok := search(c, "ou=" + group + "," + c.LDAP.DN, "(|(uid=" + uid + "))", attributes)
	if len(entries) < 1 {
		log.WithFields(log.Fields{
			"message": "ldap search successful but search resulted in 0 results",
		}).Info("ldap")
		return nil, true
	} else if !ok {
		return nil, false
	} 
	entryMap := entriesToMap(entries, attributes)
	return entryMap, true
}

func GetGroup(c *config.Config, group string) (map[string]map[string][]string, bool) {
	entries, ok := search(c, "ou=groups," + c.LDAP.DN, "(|(cn=" + group + "))", []string{"cn", "member"})
	if len(entries) < 1 {
		log.WithFields(log.Fields{
			"message": "ldap search successful but search resulted in 0 results",
		}).Error("ldap")
		return nil, true
	} else if !ok {
		return nil, false
	} 
	arr := make(map[string]map[string][]string, len(entries))
	for _, entry := range entries {
		arr[entry.GetAttributeValue("cn")] = make(map[string][]string, 1)
		for _, m := range entry.GetAttributeValues("member") {
			arr[entry.GetAttributeValue("cn")]["members"] = append(arr[entry.GetAttributeValue("cn")]["members"], strings.Split(string(m), ",")[0])
		}
	}
	return arr, true
}

func GetAllGroups(c *config.Config) (map[string]map[string][]string, bool) {
	return GetGroup(c, "*")
}