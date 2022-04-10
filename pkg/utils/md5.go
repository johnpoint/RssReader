package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5(str string) string {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(m.Sum(nil))
}
