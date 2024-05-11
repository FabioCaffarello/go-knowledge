package gomd5id

import (
	"crypto/md5"
	"encoding/hex"
)

type ID string

func GetIDFromString(data string) ID {
	hasher := md5.New()
	hasher.Write([]byte(data))
	hash := hasher.Sum(nil)
	return ID(hex.EncodeToString(hash))
}
