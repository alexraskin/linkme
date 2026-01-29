package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"net/http"

	"github.com/syumai/workers"
)

//go:embed static/index.html
var indexHTML string

//go:embed static/style.css
var styleCSS string

type Link struct {
	Title string
	URL   string
	Icon  string
}

type PageData struct {
	Name   string
	Bio    string
	Avatar string
	CSS    template.CSS
	Links  []Link
}

var config = PageData{
	Name:   "alex raskin",
	Bio:    "devops engineer • cat dad • us-west-1",
	Avatar: "https://avatars.githubusercontent.com/u/1234567", // replace with your avatar URL
	Links: []Link{
		{Title: "stand with iran", URL: "https://standwithiran.org", Icon: "❤️"},
		{Title: "cease fire today", URL: "https://ceasefiretoday.com/", Icon: "✊"},
		{Title: "know your rights", URL: "https://www.aclu.org/know-your-rights", Icon: "⚖️"},
		{Title: "cosmo the cat", URL: "https://cosmothecat.net", Icon: "🐱"},
		{Title: "github", URL: "https://github.com/alexraskin", Icon: "🐙"},
		{Title: "twitter", URL: "https://twitter.com/notalexraskin", Icon: "𝕏"},
	},
}

func main() {
	tmpl := template.Must(template.New("index").Parse(indexHTML))

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}

		data := config
		data.CSS = template.CSS(styleCSS)

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(buf.Bytes())
	})

	workers.Serve(nil)
}
