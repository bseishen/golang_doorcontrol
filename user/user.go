package user

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

type User struct {
	Active  int
	Key     int
	DBHash  string
	IrcName string
}

func (u *User) Clear() {
	u.Active = 0
	u.Key = 0
	u.IrcName = ""
	u.DBHash = ""
}

func (u *User) EncryptPass(pass string) string {
	b := sha1.Sum([]byte(pass))
	var c []byte = b[:]
	return hex.EncodeToString(c)
}

func (u *User) ValidatePass(pass string) bool {
	if strings.Compare(u.EncryptPass(pass), u.DBHash) == 0 {
		return true
	} else {
		return false
	}
}
