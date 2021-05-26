package main

import (
	"log"
	"net/http"

	"gee"
)

func main() {
	r := gee.New()

	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page.</h1>")
	})

	v1 := r.Group("/v1")
	{

		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "Hello  %s, you're at %s", c.Query("name"), c.Path)
		})

	}

	v2 := r.Group("/v2")
	{
		v2.POST("/login", func(c *gee.Context) {
			c.Json(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		v2.GET("/assets/*filepath", func(c *gee.Context) {
			c.Json(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
		})
	}
	log.Fatal(r.Run(":8199"))
}
