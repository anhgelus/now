package main

import (
	"embed"
	"encoding/json"
	"flag"
	"github.com/anhgelus/golatt"
	"log/slog"
	"os"
)

//go:embed templates
var templates embed.FS

var (
	domain   string
	dataPath string
)

func init() {
	flag.StringVar(&domain, "domain", "", "domain to use")
	flag.StringVar(&dataPath, "data", "", "data to use")
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
	if dataPath == "" {
		dataPath = os.Getenv("NOW_DATA")
		if dataPath == "" {
			slog.Error("Data not set. Set it with --data relative path or with the env NOW_DATA")
			return
		}
	}
	b, err := os.ReadFile(dataPath)
	if err != nil {
		panic(err)
	}
	var data Data
	err = json.Unmarshal(b, &data)
	if err != nil {
		panic(err)
	}
	g := golatt.New(templates)
	g.DefaultSeoData = &golatt.SeoData{
		Image:       data.Image,
		Description: data.Description,
		Domain:      domain,
	}
	g.Templates = append(g.Templates, "templates/base/*.gohtml")

	home := golatt.Template{
		Golatt: g,
		Name:   "index",
		Title:  data.Person.Name,
		Data:   &data,
		URL:    "/",
	}
	credits := golatt.Template{
		Golatt: g,
		Name:   "credits",
		Title:  "Credits",
		Data:   &data,
		URL:    "/credits",
	}

	g.HandleFunc("/", home.Handle())
	g.HandleFunc("/credits", credits.Handle())

	g.StartServer(":80")
}
