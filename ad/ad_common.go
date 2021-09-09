package ad

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/encoding/unicode"
	"strings"
	"uidm/common"
	"uidm/output"
)

func BindAD(ldapUrl, username, password string) *ldap.Conn {
	// connect to AD
	ldapConn, err := ldap.DialURL(ldapUrl, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	output.GenerateLog(err, "AD bind", ldapUrl, false)

	_, err = ldapConn.SimpleBind(&ldap.SimpleBindRequest{
		Username: username,
		Password: password,
	})
	output.GenerateLog(err, "AD Auth", username, false)

	return ldapConn
}

func GenerateADPassword(pwd string) string {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	pwdEncoded, err := utf16.NewEncoder().String(fmt.Sprintf("%q", pwd))
	output.GenerateLog(err, "AD password generation", pwd, false)

	return pwdEncoded
}

func SearchADObject(l *ldap.Conn, baseDN, n, searchAttr string, resultAttr []string) []*ldap.Entry {
	filter := fmt.Sprintf("(%s=%s)", searchAttr, ldap.EscapeFilter(n))

	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, resultAttr, []ldap.Control{})
	resultSearchUser, err := l.Search(searchReq)
	output.GenerateLog(err, "Search AD object", n, false)

	return resultSearchUser.Entries
}

func GetADUserGroupManager(l *ldap.Conn, baseDN, uSAM string) map[string]string {
	mgr := make(map[string]string)

	uEntry := SearchADObject(l, baseDN, uSAM, "sAMAccountName", []string{"memberOf"})
	gDNs := uEntry[0].GetAttributeValues("memberOf")
	if len(gDNs) > 0 {
		for _, gDN := range gDNs {
			gCN := common.ConvertDNToCN(gDN)
			gEntry := SearchADObject(l, baseDN, gCN, "cn", []string{"managedBy"})
			mgrDN := gEntry[0].GetAttributeValue("managedBy")
			if mgrDN != "" {
				mgrCN := common.ConvertDNToCN(mgrDN)
				mgrEntry := SearchADObject(l, baseDN, mgrCN, "cn", []string{"mail"})
				mgrMail := mgrEntry[0].GetAttributeValue("mail")
				if mgrMail != "" {
					mgr[gCN] = mgrMail
				}
			}
		}
	}
	return mgr
}

func CreateADUser(l *ldap.Conn, attr map[string][]string, dn string) error {
	req := ldap.NewAddRequest(dn, []ldap.Control{})

	for k, v := range attr {
		req.Attribute(k, v)

	}
	err := l.Add(req)
	output.GenerateLog(err, "AD user creation", dn, false)
	return err
}

func ModifyADUser(l *ldap.Conn, attr map[string][]string, dn string) error {
	req := ldap.NewModifyRequest(dn, []ldap.Control{})

	for k, v := range attr {
		if len(v) > 0 && v[0] != "" {
			req.Replace(k, v)
		} else {
			req.Delete(k, []string{})
		}
	}
	err := l.Modify(req)
	output.GenerateLog(err, "AD user modification", dn, false)
	return err
}

func DisableADUser(l *ldap.Conn, dn string) error {
	req := ldap.NewModifyRequest(dn, []ldap.Control{})

	req.Replace("userAccountControl", []string{fmt.Sprintf("%d", 0x0200)})
	err := l.Modify(req)
	output.GenerateLog(err, "AD user disable", dn, false)

	return err
}

func AddADGroupMember(l *ldap.Conn, baseDN, grp, usr string) {
	usrSAM := common.ConvertDNToCN(usr)
	e := SearchADObject(l, baseDN, usrSAM, "cn", []string{"memberOf"})

	var isMember bool
	if len(e) > 0 {
		for _, v := range e[0].GetAttributeValues("memberOf") {
			if grp == v {
				isMember = true
				break
			}
		}
	}
	// user is not a member yet
	if !isMember {
		req := ldap.NewModifyRequest(grp, []ldap.Control{})
		req.Add("member", []string{usr})

		err := l.Modify(req)
		output.GenerateLog(err, "Adding AD user to group", usr+" | "+grp, false)
	}
}

func RemoveADGroupMember(l *ldap.Conn, baseDN, grp, usr string) {
	usrSAM := common.ConvertDNToCN(usr)
	e := SearchADObject(l, baseDN, usrSAM, "cn", []string{"memberOf"})

	var isMember bool
	if len(e) > 0 {
		for _, v := range e[0].GetAttributeValues("memberOf") {
			if grp == v {
				isMember = true
				break
			}
		}
	}
	// user is a member
	if isMember {
		req := ldap.NewModifyRequest(grp, []ldap.Control{})
		req.Delete("member", []string{usr})

		err := l.Modify(req)
		output.GenerateLog(err, "Removing AD user from group", usr+" | "+grp, false)
	}
}

func MoveADObject(l *ldap.Conn, oldDN, newDN string) {
	destOU := newDN[strings.Index(newDN, ",")+1:]
	req := ldap.NewModifyDNRequest(oldDN, "CN="+common.ConvertDNToCN(oldDN), true, destOU)

	err := l.ModifyDN(req)
	output.GenerateLog(err, "Moving AD object", oldDN+" | "+newDN, false)
}

func ConvertDNToSAM(l *ldap.Conn, baseDN, dn string) string {
	cn := common.ConvertDNToCN(dn)
	e := SearchADObject(l, baseDN, cn, "cn", []string{"sAMAccountName"})
	return e[0].GetAttributeValue("sAMAccountName")
}
