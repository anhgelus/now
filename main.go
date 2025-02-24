package main

import (
	"embed"
	"encoding/json"
	"flag"
	"github.com/anhgelus/golatt"
	"log/slog"
	"net/http"
	"os"
)

var (
	//go:embed templates
	templates embed.FS
	//go:embed dist
	assets embed.FS
)

var (
	domain     string
	configPath string
	dev        bool
)

func init() {
	flag.StringVar(&domain, "domain", "", "domain to use")
	flag.StringVar(&configPath, "config", "", "config to use")
	flag.BoolVar(&dev, "dev", false, "dev mode enabled")
}

func main() {
	flag.Parse()
	if domain == "" {
		domain = os.Getenv("NOW_DOMAIN")
		if domain == "" {
			slog.Error("Domain not set. Set it with --domain value or with the env NOW_DOMAIN")
			return
		}
	}
	if configPath == "" {
		configPath = os.Getenv("NOW_DATA")
		if configPath == "" {
			slog.Error("Config not set. Set it with --config relative path or with the env NOW_DATA")
			return
		}
	}
	b, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	var cfg Config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		panic(err)
	}
	customPages, err := cfg.LoadCustomPages()
	if err != nil {
		panic(err)
	}
	var g *golatt.Golatt
	if dev {
		g = golatt.New(templates, os.DirFS("public"), os.DirFS("dist"))
	} else {
		g = golatt.New(templates, os.DirFS("public"), golatt.UsableEmbedFS("dist", assets))
	}
	g.DefaultSeoData = &golatt.SeoData{
		Image:       cfg.Image,
		Description: cfg.Description,
		Domain:      domain,
	}
	g.Templates = append(g.Templates, "templates/base/*.gohtml")

	g.NewTemplate("index", "/", cfg.Person.Name, "", "", &cfg).Handle()
	g.NewTemplate("legal",
		"/legal",
		"Legal things",
		"",
		"Legal information about "+cfg.Person.Name+"'s Now page",
		&cfg).
		Handle()
	g.NewTemplate("tags",
		"/tags",
		"Tags",
		"",
		"Tags of "+cfg.Person.Name+"'s Now page",
		&cfg).
		Handle()

	for _, cp := range customPages {
		slog.Info("Creating custom page...", "title", cp.Title, "uri", cp.URI)
		g.NewTemplate("custom_page",
			cp.URI,
			cp.Title,
			cp.Image,
			cp.Description,
			cp).
			Handle()
	}

	g.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	if dev {
		slog.Info("Starting on http://localhost:8000/")
		g.StartServer(":8000")
	} else {
		g.StartServer(":80")
	}
}
