package main

import (
	"WebFramework/web"
	"fmt"
)

func main() {
	s := web.NewHttpServer()
	fmt.Println(s.MatchRoute("/user/:id(^[0-9]+$)", "/user/123"))
	fmt.Println(s.MatchRoute("/user/:id(^[0-9]+$)", "/user/qweqrw"))

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

	//r := gin.Default()
	//r.GET("/abc/:a", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"a":      c.Param("a"),
	//		"method": "/abc/:a",
	//	})
	//})
	//r.GET("/abc/:a/:b", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"a":      c.Param("a"),
	//		"b":      c.Param("b"),
	//		"method": "/abc/:a/:b",
	//	})
	//})
	//
	//_ = r.Run(":8080")

	//compile, _ := regexp.Compile("^[0-9]+$")
	//r, _ := regexp.Compile("^[0-9]+$")
	//fmt.Println(reflect.DeepEqual(compile, r))
	//fmt.Println(compile.String())

	//less := Less(1, 2)
	//fmt.Println(less)
}

func Less[K int | float32](a, b K) bool {
	return a < b
}
