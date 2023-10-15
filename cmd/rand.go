package main

import (
	utils "Sentinel/utils/string"
	"log"
)

func main() {
	d := utils.RandString(64)

	log.Println(d)
}
