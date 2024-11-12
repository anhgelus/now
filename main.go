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

//go:embed templates
var templates embed.FS

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
			slog.Error("Config not set. Set it with --cfg relative path or with the env NOW_DATA")
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
	g := golatt.New(templates)
	g.DefaultSeoData = &golatt.SeoData{
		Image:       cfg.Image,
		Description: cfg.Description,
		Domain:      domain,
	}
	g.Templates = append(g.Templates, "templates/base/*.gohtml")

	g.NewTemplate("index", "/", cfg.Person.Name, "", "", &cfg).Handle()
	g.NewTemplate("credits",
		"/credits",
		"Credits",
		"",
		"Credits of "+cfg.Person.Name+"'s Now page",
		&cfg).
		Handle()
	g.NewTemplate("tags",
		"/tags",
		"Tags",
		"",
		"Tags of "+cfg.Person.Name+"'s Now page",
		&cfg).
		Handle()

	g.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	if dev {
		g.StartServer(":8000")
	} else {
		g.StartServer(":80")
	}
}
