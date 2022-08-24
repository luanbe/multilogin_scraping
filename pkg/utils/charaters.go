package util

import (
	"log"
	"regexp"
	"strings"
)

func RemoveSpecialCharacters(str string) string {
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	str = re.ReplaceAllString(str, "")
	str = strings.Replace(str, "None", "", -1)
	return str
}
