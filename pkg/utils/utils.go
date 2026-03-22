package utils

import (
	"log"
	"strings"
)

func CreateId(uid, toUid string) string {
	return uid + "->" + toUid
}
func GetId(id string) (string, string) {
	log.Println(id)
	parts := strings.Split(id, "->")
	log.Println("begin")
	if len(parts) == 2 {
		log.Println("part[0]", parts[0], "part[1]", parts[1])
		return parts[0], parts[1]
	}
	log.Println("false")
	return "", ""
}
