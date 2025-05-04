package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// GetMD5 calculates the MD5 hash of a given string
func GetMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
