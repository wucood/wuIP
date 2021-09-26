package main

import "github.com/gin-gonic/gin"

func main() {
	r:=gin.Default()
	r.GET("",indexHandler)
	r.GET("/api/:ip",ipGeoLiteHandler)
	r.Run()
}
