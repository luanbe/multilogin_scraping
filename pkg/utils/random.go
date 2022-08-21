package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RandIntRange(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandInt() int {
	return rand.Int()
}

func RandMapStr(strMap map[string]string) (string, string) {
	i := rand.Intn(len(strMap))
	for k := range strMap {
		if i == 0 {
			return k, strMap[k]
		}
		i--
	}
	panic("Never")
}

func RandSliceStr(strSlice []string) string {
	return strSlice[rand.Intn(len(strSlice))]
}
