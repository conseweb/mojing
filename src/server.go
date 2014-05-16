package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	// "html/template"
	// "log"
	// "net/http"
)

func main() {
	m := martini.Classic()
	// m.Get("/", func() string {
	// 	return "Hello world!"
	// })

	// render html templates from templates directory
	m.Use(render.Renderer(render.Options{
		Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
	}))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "hello", "魔镜")
	})

	// This will set the Content-Type header to "application/json; charset=UTF-8"
	m.Get("/api", func(r render.Render) {
		r.JSON(200, map[string]interface{}{"hello": "world"})
	})

	m.Run()
	// log.Fatal(http.ListenAndServe(":3000", m))
}
