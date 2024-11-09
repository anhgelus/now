package main

import (
	"html/template"
	"strconv"
)

type Data struct {
	Image           string           `json:"image"`
	Description     string           `json:"description"`
	Person          *Person          `json:"person"`
	BackgroundColor *BackgroundColor `json:"background_color"`
	Links           []*Link          `json:"links"`
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

type Link struct {
	Link    string `json:"link"`
	Content string `json:"content"`
	Color   string `json:"color"`
}

func (d *Data) GetBackground() template.CSS {
	bg := d.BackgroundColor
	css := "background: " + bg.Type + "-gradient("
	if bg.Type == "linear" {
		css += strconv.Itoa(int(bg.Angle)) + "deg,"
	}
	for _, c := range bg.Colors {
		css += c.Color + " " + strconv.Itoa(int(c.Position)) + "%,"
	}
	return template.CSS(css[:len(css)-1] + ");")
}

func (l *Link) GetBackground() template.CSS {
	return template.CSS("background: " + l.Color)
}
