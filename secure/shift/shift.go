// Package shift
//
// @author: xwc1125
package shift

// Shift 位移变换key
func Shift(key string, tag int) string {
	keyLen := len(key)
	if keyLen == 0 {
		return ""
	}
	result := make([]byte, keyLen)
	tempKeyBytes := []byte(key)
	var k = 0
	removeIndexs := make([]int, 0)

	for i := tag + 1; i > 0; i-- {
		size := len(tempKeyBytes)
		removeIndexs = make([]int, 0)
		for j := 0; j < size; j++ {
			if j%i == 0 {
				if k < keyLen {
					aa := tempKeyBytes[j]
					result[k] = aa
					k++
					removeIndexs = append(removeIndexs, j)
				} else {
					break
				}
			}
		}
		for l := 0; l < len(removeIndexs); l++ {
			index := removeIndexs[l] - l
			if index == 0 {
				tempKeyBytes = tempKeyBytes[index+1:]
			} else if index == len(tempKeyBytes)-1 {
				tempKeyBytes = tempKeyBytes[:index]
			} else {
				tempKeyBytes = append(tempKeyBytes[0:index], tempKeyBytes[index+1:]...)
			}
		}
	}
	return string(result)
}
