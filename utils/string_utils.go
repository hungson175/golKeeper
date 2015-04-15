package utils

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"runtime"
)

func GetMD5Hash(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
