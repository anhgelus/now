package main

import (
	"github.com/anhgelus/golatt"
	"html/template"
	"strconv"
)

type Config struct {
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Person      *Person `json:"person"`
	Color       *Color  `json:"colors"`
	Links       []*Link `json:"links"`
	Legal       *Legal  `json:"legal"`
}

type Person struct {
	Name     string `json:"name"`
	Pronouns string `json:"pronouns"`
	Image    string `json:"image"`
	Tags     []*Tag `json:"tags"`
}

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Color struct {
	Background *BackgroundColor `json:"background"`
	Text       string           `json:"text"`
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
	Link           string `json:"link"`
	Content        string `json:"content"`
	Color          string `json:"color"`
	TextColor      string `json:"text_color"`
	ColorHover     string `json:"color_hover"`
	TextColorHover string `json:"text_color_hover"`
}

type Legal struct {
	LegalInformationLink string   `json:"legal_information_link"`
	ImagesSource         []string `json:"images_source"`
	FontSource           string   `json:"font_source"`
}

func (c *Config) GetBackground() template.CSS {
	bg := c.Color.Background
	css := "background: " + bg.Type + "-gradient("
	if bg.Type == "linear" {
		css += strconv.Itoa(int(bg.Angle)) + "deg,"
	}
	for _, c := range bg.Colors {
		css += c.Color + " " + strconv.Itoa(int(c.Position)) + "%,"
	}
	return template.CSS(css[:len(css)-1] + ");")
}

func (c *Config) GetBackgroundImage() template.CSS {
	return template.CSS("--background-image: url(" + golatt.GetStaticPath(c.Image) + ");")
}

func (c *Config) GetTextColor() template.CSS {
	return template.CSS("--text-color: " + c.Color.Text + ";")
}

func (l *Link) GetLinkColor() template.CSS {
	return template.CSS("--text-color: " + l.TextColor + ";--text-color-hover: " + l.TextColorHover + ";")
}

func (l *Link) GetBackground() template.CSS {
	return template.CSS("--background: " + l.Color + ";--background-hover: " + l.ColorHover + ";")
}
