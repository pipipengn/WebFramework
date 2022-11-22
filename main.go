package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {

	//s := web.NewHttpServer(web.WithMiddleware(), web.WithMiddleware())

	//s.Post("/", func(c *web.Context) {
	//	data := Form{}
	//	_ = c.BindJSON(&data)
	//	c.JSON(200, data)
	//	c.PathValue("qwe")
	//})
	//_ = s.Start(":8080")

	//fmt.Println(s.MatchRoute("/user/:id(^[0-9]+$)", "/user/123"))
	//fmt.Println(s.MatchRoute("/user/:id(^[0-9]+$)", "/user/qweqrw"))

	//s.Get("/user/:id(^[0-9]+$)", func(c *web.Context) {
	//	c.JSON("qwer")
	//})
	//
	//s.Get("/order/:pid/:mid", func(c *web.Context) {
	//	m := map[string]string{"pid": c.Param("pid"), "mid": c.Param("mid")}
	//	c.JSON(m)
	//})
	//
	//s.Get("/abc/:a", func(c *web.Context) {
	//	c.JSON(map[string]string{
	//		"a":      c.Param("a"),
	//		"method": "/abc/:a",
	//	})
	//})
	//
	//s.Get("/abc/:a/:b", func(c *web.Context) {
	//	c.JSON(map[string]string{
	//		"a":      c.Param("a"),
	//		"b":      c.Param("b"),
	//		"method": "/abc/:a/:b",
	//	})
	//})

	//_ = s.Start(":8080")

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		fmt.Println(123)
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		//data := Form{}
		//_ = c.BindJSON(&data)
		//c.JSON(200, data)

		fmt.Println("qweqweqweqwe")
		c.JSON(200, gin.H{})
	})

	_ = r.Run(":8080")
}
