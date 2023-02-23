// Package ddm
//
// @author: xwc1125
package ldap

import (
	"crypto/tls"
	"fmt"

	"github.com/chain5j/logger"
	"github.com/go-ldap/ldap/v3"
)

type LdapConfig struct {
	Enable          bool           `json:"enable" mapstructure:"enable"`
	Host            string         `json:"host" mapstructure:"host"`
	Port            int            `json:"port" mapstructure:"port"`
	BaseDn          string         `json:"base_dn" mapstructure:"base_dn"`
	BindUser        string         `json:"bind_user" mapstructure:"bind_user"`
	BindPass        string         `json:"bind_pass" mapstructure:"bind_pass"`
	AuthFilter      string         `json:"auth_filter" mapstructure:"auth_filter"`
	Attributes      ldapAttributes `json:"attributes" mapstructure:"attributes"`
	CoverAttributes bool           `json:"cover_attributes" mapstructure:"cover_attributes"`
	TLS             bool           `json:"tls" mapstructure:"tls"`
	StartTLS        bool           `json:"start_tls" mapstructure:"start_tls"`
}

type ldapAttributes struct {
	Nickname string `json:"nickname" mapstructure:"nickname"`
	Phone    string `json:"phone" mapstructure:"phone"`
	Email    string `json:"email" mapstructure:"email"`
	UID      string `json:"uid" mapstructure:"uid"`
}

var LDAP LdapConfig

func InitLdap(ldap LdapConfig) {
	LDAP = ldap
}

func genLdapAttributeSearchList() []string {
	var ldapAttributes []string
	attrs := LDAP.Attributes
	if attrs.Nickname != "" {
		ldapAttributes = append(ldapAttributes, attrs.Nickname)
	}
	if attrs.Email != "" {
		ldapAttributes = append(ldapAttributes, attrs.Email)
	}
	if attrs.Phone != "" {
		ldapAttributes = append(ldapAttributes, attrs.Phone)
	}
	if attrs.UID != "" {
		ldapAttributes = append(ldapAttributes, attrs.UID)
	}
	return ldapAttributes
}

func LdapReq(user, pass string) (*ldap.SearchResult, error) {
	var conn *ldap.Conn
	var err error
	lc := LDAP
	addr := fmt.Sprintf("%s:%d", lc.Host, lc.Port)

	if lc.TLS {
		conn, err = ldap.DialTLS("tcp", addr, &tls.Config{InsecureSkipVerify: true})
	} else {
		conn, err = ldap.Dial("tcp", addr)
	}

	if err != nil {
		logger.Error("cannot dial ldap", "addr", addr, "err", err)
		return nil, internalServerError
	}

	defer conn.Close()

	if !lc.TLS && lc.StartTLS {
		if err := conn.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			logger.Error("conn startTLS fail", "err", err)
			return nil, internalServerError
		}
	}

	// if bindUser is empty, anonymousSearch mode
	if lc.BindUser != "" {
		// BindSearch mode
		if err := conn.Bind(lc.BindUser, lc.BindPass); err != nil {
			logger.Error("bind ldap fail", "user", lc.BindUser, "err", err)
			return nil, internalServerError
		}
	}

	searchRequest := ldap.NewSearchRequest(
		lc.BaseDn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.AuthFilter, user), // The filter to apply
		genLdapAttributeSearchList(),     // A list attributes to retrieve
		nil,
	)

	sr, err := conn.Search(searchRequest)

	if err != nil {
		logger.Error("ldap search fail", "err", err)
		return nil, internalServerError
	}

	if len(sr.Entries) == 0 {
		logger.Info("ldap auth fail, no such user", "user", user)
		return nil, loginFailError
	}

	if len(sr.Entries) > 1 {
		logger.Error("search user, multi entries found", "user", user)
		return nil, internalServerError
	}

	if err := conn.Bind(sr.Entries[0].DN, pass); err != nil {
		logger.Info("password error", "user", user, "err", err)
		return nil, loginFailError
	}
	return sr, nil
}
