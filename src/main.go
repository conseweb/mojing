package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"io"
	"sort"
	// "html/template"
	// "log"
	"net/http"
)

var TOKEN string = "mojing"

// 对字符串进行SHA1哈希
func sha1hash(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

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

	m.Get("/wx/**", func(req *http.Request) string {
		// test url:
		// http://localhost:3000/wx?signature=6dfc072a5cb27abfd8fdc5f16ada6ed34380ddff&timestamp=3436&nonce=23435&token=mojing&echostr=testok
		// get 3 params
		req.ParseForm()
		// if len(req.Form) > 0 {
		// 	for k, v := range req.Form {
		// 		x := fmt.Sprintf("%s, %s", k, v[0])
		// 		fmt.Println(x)
		// 	}
		// }
		//
		// res.Write("helloxx")
		// return "Helloxx"
		sig := req.Form.Get("signature")
		ts := req.Form.Get("timestamp")
		nonce := req.Form.Get("nonce")
		echostr := req.Form.Get("echostr")

		// fmt.Printf("1: %s, 2: %s, 3: %s\n", TOKEN, ts, nonce)
		// sort
		sa := []string{}
		sa = append(sa, TOKEN, ts, nonce)
		sort.Sort(sort.StringSlice(sa))
		// hash
		data := fmt.Sprintf("%s%s%s", sa[0], sa[1], sa[2])
		sekret := sha1hash(data)
		if sekret == sig {
			return echostr
		} else {
			return "not wechat message"
		}
	})

	m.Run()
	// log.Fatal(http.ListenAndServe(":3000", m))
}
