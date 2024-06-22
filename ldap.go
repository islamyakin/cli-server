package main

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

const (
	ldapServer   = ""
	ldapPort     =
	ldapBindDN   = ""
	ldapPassword = ""
	ldapSearchDN = ""
)

type UserLDAPData struct {
	ID       string
	Email    string
	Name     string
	FullName string
	Wifi     string
}

func AuthUsingLDAP(username, password string) (bool, *UserLDAPData, error) {
	ldapUrl := fmt.Sprintf("ldap://%s:%d", ldapServer, ldapPort)
	l, err := ldap.DialURL(ldapUrl)
	if err != nil {
		fmt.Printf("Cannot connect to ldap %s\n", err)
	}

	defer func(l *ldap.Conn) {
		err := l.Close()
		if err != nil {
			fmt.Printf("err")
		}
	}(l)

	err = l.Bind(ldapBindDN, ldapPassword)
	if err != nil {
		return false, nil, err
	}
	searchRequest := ldap.NewSearchRequest(
		ldapSearchDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username),
		[]string{"dn", "cn", "sn", "mail"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return false, nil, err
	}

	if len(sr.Entries) == 0 {
		return false, nil, fmt.Errorf("user not found")
	}
	entry := sr.Entries[0]

	err = l.Bind(entry.DN, password)
	if err != nil {
		return false, nil, err
	}
	data := new(UserLDAPData)
	data.ID = username

	for _, attr := range entry.Attributes {
		switch attr.Name {
		case "sn":
			data.Name = attr.Values[0]
		case "mail":
			data.Email = attr.Values[0]
		case "cn":
			data.FullName = attr.Values[0]

		}
	}

	return true, data, nil
}
