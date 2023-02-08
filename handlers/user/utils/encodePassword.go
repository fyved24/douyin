package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func EncodePassword(password string) string {
	d := []byte(password)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}
