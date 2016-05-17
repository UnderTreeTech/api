package util

import (
	"crypto/sha1"
	"io"
	"fmt"
)


//获取SHA1加密串
func SHA1Encrypt(str string) (encryptStr string) {
	h := sha1.New()
	io.WriteString(h, str)

	return fmt.Sprintf("%x", h.Sum(nil))
}