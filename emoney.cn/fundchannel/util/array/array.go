package _array

import (
	"math/rand"
	"time"
)

func Distinct(source []int) (target []int) {
	if source == nil {
		return nil
	}
	distinctMap := make(map[int]int)
	for _, item := range source {
		_, exist := distinctMap[item]
		if !exist {
			target = append(target, item)
			distinctMap[item] = 1
		}
	}
	return target
}

// GetShuffleArray 生成包含count个元素的打乱后的数组 数组元素为[0,count)中的数字
func GetShuffleArray(count int) (shuffleArr []int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffleArr = r.Perm(count)
	return shuffleArr
}
