package main

import (
	"embed"
	"encoding/json"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/anhgelus/golatt"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

var (
	//go:embed templates
	templates embed.FS
	//go:embed dist
	assets embed.FS
)

var (
	domain       string
	configPath   string
	dev          bool
	generateToml bool
	generateJson bool
)

func init() {
	flag.StringVar(&domain, "domain", "", "domain to use")
	flag.StringVar(&configPath, "config", "", "config to use")
	flag.BoolVar(&dev, "dev", false, "dev mode enabled")
	flag.BoolVar(&generateJson, "generate-json-config", false, "generate a config example")
	flag.BoolVar(&generateToml, "generate-toml-config", false, "generate a config example")
}

func main() {
	flag.Parse()
	if generateToml || generateJson {
		generateConfigFile(generateToml)
		return
	}
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
	if strings.HasSuffix(configPath, ".json") {
		err = json.Unmarshal(b, &cfg)
	} else if strings.HasSuffix(configPath, ".toml") {
		err = toml.Unmarshal(b, &cfg)
	} else {
		panic("config file must be .json or .toml")
	}
	if err != nil {
		panic(err)
	}
	customPages, err := cfg.LoadCustomPages()
	if err != nil {
		panic(err)
	}

	var g *golatt.Golatt
	if dev {
		g = golatt.New(golatt.UsableEmbedFS("templates", templates), os.DirFS("public"), os.DirFS("dist"))
	} else {
		g = golatt.New(
			golatt.UsableEmbedFS("templates", templates),
			os.DirFS("public"),
			golatt.UsableEmbedFS("dist", assets),
		)
	}
	g.DefaultSeoData = &golatt.SeoData{
		Image:       cfg.Image,
		Description: cfg.Description,
		Domain:      domain,
	}
	g.Templates = append(g.Templates, "base/*.gohtml")

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

func generateConfigFile(isToml bool) {
	cfg := Config{
		Image:       "wallpaper.webp",
		Description: "I am a beautiful description!",
		Person: &Person{
			Name:     "John Doe",
			Pronouns: "any",
			Image:    "pfp.webp",
			Tags: []*Tag{
				{"Hello", "World", ""},
				{"I am", "a tag", ""},
			},
		},
		Color: &Color{
			Background: &BackgroundColor{
				Type:  "linear",
				Angle: 141,
				Colors: []struct {
					Color    string `json:"color" toml:"color"`
					Position uint   `json:"position" toml:"position"`
				}{
					{"#a4a2b8", 0},
					{"#3b3860", 40},
					{"#0f0c2c", 80},
				},
			},
			Button: &ButtonColor{
				Text:            "#4c0850",
				TextHover:       "#57145b",
				Background:      "#f399d0",
				BackgroundHover: "#f5c0e0",
			},
			Text:     "#fff",
			TagHover: "#000",
		},
		Links: []*Link{
			{"/foo", "Blog"},
			{"https://www.youtube.com/@anhgelus", "YouTube"},
		},
		Legal: &Legal{
			LegalInformationLink: "/bar",
			ImagesSource: []string{
				"Profile picture: some one on website",
				"Background picture: another one on another website",
			},
			FontSource: "Name by some one on website",
		},
		CustomPages: []string{"custom.json"},
	}
	var b []byte
	var err error
	if isToml {
		b, err = toml.Marshal(&cfg)
	} else {
		b, err = json.Marshal(&cfg)
	}
	if err != nil {
		panic(err)
	}
	println(string(b))
}
