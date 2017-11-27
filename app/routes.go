package main

func initializeRoutes() {
	router.GET("/", showIndexPage)
	router.GET("/ws/", wsConnection)
}
