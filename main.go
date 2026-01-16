package main

import (
	"embed"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var configFile []byte

//go:embed templates/*
var templates embed.FS

//go:embed static/*
var static embed.FS

//go:embed all:assets
var assets embed.FS

type Config struct {
	ConfigVersion string     `yaml:"configVersion"`
	Meta          Meta       `yaml:"meta"`
	Name          string     `yaml:"name"`
	Subtitle      string     `yaml:"subtitle"`
	Description   string     `yaml:"description"`
	Avatar        string     `yaml:"avatar"`
	Theme         string     `yaml:"theme"`
	Background    Background `yaml:"background"`
	Links         []Link     `yaml:"links"`
	Sections      []Section  `yaml:"sections"`
	Socials       []Social   `yaml:"socials"`
	Footer        Footer     `yaml:"footer"`
}

type Meta struct {
	SiteURL     string `yaml:"siteUrl"`
	Description string `yaml:"description"`
	Favicon     string `yaml:"favicon"`
}

type Background struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

type Link struct {
	Title   string `yaml:"title"`
	URL     string `yaml:"url"`
	Icon    string `yaml:"icon"`
	IconURL string `yaml:"iconUrl"`
	Color   string `yaml:"color"`
}

type Section struct {
	Title string `yaml:"title"`
	Links []Link `yaml:"links"`
}

type Social struct {
	Icon  string `yaml:"icon"`
	URL   string `yaml:"url"`
	Color string `yaml:"color"`
}

type Footer struct {
	PoweredBy PoweredBy `yaml:"poweredBy"`
}

type PoweredBy struct {
	Text string `yaml:"text"`
	URL  string `yaml:"url"`
}

type PageData struct {
	Config
	BuildTime string
}

var config Config

var buildTime = "unknown"

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	if err := yaml.Unmarshal(configFile, &config); err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5))

	funcMap := template.FuncMap{
		"trimSpace": strings.TrimSpace,
		"splitLines": func(s string) []string {
			lines := strings.Split(strings.TrimSpace(s), "\n")
			var result []string
			for _, l := range lines {
				if t := strings.TrimSpace(l); t != "" {
					result = append(result, t)
				}
			}
			return result
		},
	}

	tmpl := template.Must(template.New("").Funcs(funcMap).ParseFS(templates, "templates/index.html"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data := PageData{
			Config:    config,
			BuildTime: buildTime,
		}
		tmpl.ExecuteTemplate(w, "index.html", data)
	})

	r.Handle("/static/*", http.FileServer(http.FS(static)))
	r.Handle("/assets/*", http.FileServer(http.FS(assets)))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})

	slog.Info("Starting server", slog.String("listen_addr", ":"+strconv.Itoa(*port)))
	if err := http.ListenAndServe(":"+strconv.Itoa(*port), r); err != nil {
		panic(err)
	}
}
