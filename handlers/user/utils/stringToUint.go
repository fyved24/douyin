package utils

import "strconv"

// 将字符串转为无符号整数(转换userID用)
func StringToUint(s string) uint {
	sInt, _ := strconv.Atoi(s)
	sUint := uint(sInt)
	return sUint
}
