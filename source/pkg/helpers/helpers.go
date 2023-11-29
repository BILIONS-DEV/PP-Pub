package helpers

import (
	"encoding/json"
	"log"
)

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func NewPrint(i interface{}) {
	log.Println(PrettyPrint(i))
}
