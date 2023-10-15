package main

import (
	"Sentinel/dao"
	"Sentinel/router"
)

func main() {
	dao.InitDatabase()

	router.InitRouter()
}
