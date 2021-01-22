package main

import "Lottery/initRouter"

func main() {
	router := initRouter.SetupRouter()
	_ = router.Run()
}