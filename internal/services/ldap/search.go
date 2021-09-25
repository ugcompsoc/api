package ldap

import (
	log "github.com/sirupsen/logrus"
	"github.com/go-ldap/ldap/v3"
	"strings"
	"strconv"
)

/*
 * LDAP Utils
 */

 func entriesToMap(entries []*ldap.Entry) (map[string]*User) {
	users := make(map[string]*User, len(entries))
	for _, entry := range entries {
		uid := entry.GetAttributeValue("uid")
		users[uid] = entryToUser(entry)
	}
	return users
}

func entryToUser(entry *ldap.Entry) *User {
	uidNumber, _ := strconv.Atoi(entry.GetAttributeValue("uidNumber"))
	return &User{
		FullName: entry.GetAttributeValue("cn"),
		FirstName: entry.GetAttributeValue("givenName"),
		Surname: entry.GetAttributeValue("sn"),
		MemberID: entry.GetAttributeValue("employeeNumber"),
		Mail: entry.GetAttributeValues("mail"),
		UID: entry.GetAttributeValue("uid"),
		ObjectClass: entry.GetAttributeValues("objectClass"),
		Shell: entry.GetAttributeValue("loginShell"),
		Home: entry.GetAttributeValue("homeDirectory"),
		UIDNumber: uidNumber,
	}
}
 
/*
 * Account Search Utils
 */

func (c *Client) search(searchBase string, filter string, attributes ...string) ([]*ldap.Entry, bool) {
	if len(attributes) < 1 {
		attributes = c.Attributes
	}
	
	l := c.bind()
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

func (c *Client) GetUsersFromOU(ou string) (map[string]*User, bool) {
	entries, ok := c.search("ou=" + ou + "," + c.DN, "(|(uid=*))")
	if len(entries) < 1 {
		log.WithFields(log.Fields{
			"message": "ldap search successful but search resulted in 0 results",
		}).Info("ldap")
		return nil, true
	} else if !ok {
		return nil, false
	} 
	return entriesToMap(entries), true
}

func (c *Client) GetUser(uid string) (*User, bool) {
	entries, ok := c.search(c.DN, "(|(uid=" + uid + "))")
	if len(entries) < 1 {
		log.WithFields(log.Fields{
			"message": "ldap search successful but search resulted in 0 results",
		}).Info("ldap")
		return nil, true
	} else if !ok {
		return nil, false
	} 
	return entryToUser(entries[0]), true
}

func (c *Client) GetGroup(group string) (map[string]map[string][]string, bool) {
	entries, ok := c.search("ou=groups," + c.DN, "(|(cn=" + group + "))")
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

func (c *Client) GetAllGroups() (map[string]map[string][]string, bool) {
	return c.GetGroup("*")
}