package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/syumai/workers"
)

//go:embed templates/index.html
var indexHTML string

//go:embed static/style.css
var styleCSS string

type Link struct {
	Title string
	URL   string
	Icon  string
}

type About struct {
	Name    string
	Bio     string
	Avatar  string
	Favicon string
	CSS     template.CSS
	Links   []Link
}

var about = About{
	Name:    "alex raskin",
	Bio:     "devops engineer • cat dad • us-west-1",
	Avatar:  "https://cdn.raskin.io/avatar.jpg",
	Favicon: "https://cdn.raskin.io/favicon.ico",
	Links: []Link{
		{Title: "stand with iran", URL: "https://standwithiran.org", Icon: "❤️"},
		{Title: "cease fire today", URL: "https://ceasefiretoday.com/", Icon: "✊"},
		{Title: "know your rights", URL: "https://www.aclu.org/know-your-rights", Icon: "⚖️"},
		{Title: "cosmo the cat", URL: "https://cosmothecat.net", Icon: "🐱"},
		{Title: "github", URL: "https://github.com/alexraskin", Icon: "🐙"},
		{Title: "twitter", URL: "https://twitter.com/notalexraskin", Icon: "𝕏"},
		{Title: "bluesky", URL: "https://bsky.app/profile/alexraskin.bsky.social", Icon: "🐳"},
		{Title: "last.fm", URL: "https://www.last.fm/user/alexraskin", Icon: "🎧"},
	},
}

func main() {
	tmpl := template.Must(template.New("index").Parse(indexHTML))
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		data := about
		data.CSS = template.CSS(styleCSS)

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(buf.Bytes())
	})

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/", http.StatusMovedPermanently)
	})

	workers.Serve(r)
}
