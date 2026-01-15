package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
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

// Config represents the YAML configuration
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

// PageData combines config with runtime data
type PageData struct {
	Config
	BuildTime string
}

var config Config

// Set at build time via -ldflags
var buildTime = "unknown"

func main() {
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
