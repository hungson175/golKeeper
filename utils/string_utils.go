package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Hash(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}
