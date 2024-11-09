package main

import (
	"embed"
	"flag"
	"github.com/anhgelus/golatt"
	"log/slog"
)

//go:embed templates
var templates embed.FS

var (
	domain string
	data   string
)

func init() {
	flag.StringVar(&domain, "domain", "", "domain to use")
	flag.StringVar(&data, "data", "", "data to use")
}

func main() {
	flag.Parse()
	if domain == "" {
		slog.Error("Domain not set. Set it with --domain value")
		return
	}
	if data == "" {
		slog.Error("Data not set. Set it with --data relative path")
		return
	}
	g := golatt.New(templates)
	g.DefaultSeoData = &golatt.SeoData{
		Image:       "",
		Description: "",
		Domain:      domain,
	}
	g.Templates = append(g.Templates, "templates/page/*.gohtml")

	//g.StartServer(":80")
}
