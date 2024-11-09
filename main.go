package main

import (
	"embed"
	"github.com/anhgelus/golatt"
)

//go:embed templates
var templates embed.FS

func main() {
	g := golatt.New(templates)
	g.DefaultSeoData = &golatt.SeoData{
		Image:       "",
		Description: "",
		Domain:      "now.anhgelus.world",
	}
	g.Templates = append(g.Templates, "templates/page/*.gohtml")

	g.StartServer(":80")
}
