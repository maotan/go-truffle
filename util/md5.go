/**
* @Author: mo tan
* @Description:
* @Date 2021/1/10 12:38
 */
package util

import (
	"crypto/md5"
	"fmt"
)

func GenMd5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}
