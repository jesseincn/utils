package encrypt

import (
	"crypto/md5"
	"fmt"
)

func MD5(encode string) string {
	data := []byte(encode)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str1
}
