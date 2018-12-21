package _crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(str string) (hash string) {
	md5Provider := md5.New()
	md5Provider.Write([]byte(str))
	sum := md5Provider.Sum(nil)
	hash = hex.EncodeToString(sum)
	return hash
}
