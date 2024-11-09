package main

import (
	"html/template"
	"log/slog"
	"strconv"
)

type Data struct {
	Image           string           `json:"image"`
	Description     string           `json:"description"`
	Person          *Person          `json:"person"`
	BackgroundColor *BackgroundColor `json:"background_color"`
}

type Person struct {
	Name     string `json:"name"`
	Pronouns string `json:"pronouns"`
	Image    string `json:"image"`
}

type BackgroundColor struct {
	Type   string `json:"type"`
	Angle  uint   `json:"angle"`
	Colors []struct {
		Color    string `json:"color"`
		Position uint   `json:"position"`
	} `json:"colors"`
}

func (d *Data) GetBackground() template.CSS {
	bg := d.BackgroundColor
	css := "background: " + bg.Type + "-gradient("
	slog.Info(css)
	if bg.Type == "linear" {
		css += strconv.Itoa(int(bg.Angle)) + "deg,"
	}
	slog.Info(css)
	for _, c := range bg.Colors {
		css += c.Color + " " + strconv.Itoa(int(c.Position)) + "%,"
		slog.Info(css)
	}
	slog.Info(css[:len(css)-1] + ");")
	return template.CSS(css[:len(css)-1] + ");")
}
