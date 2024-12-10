package main

import (
	"encoding/json"
	"github.com/anhgelus/golatt"
	"html/template"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const (
	TitleContentType       = "title"
	SubtitleContentType    = "subtitle"
	ParagraphContentType   = "paragraph"
	ListContentType        = "list"
	OrderedListContentType = "ordered_list"
)

type Config struct {
	Image       string   `json:"image"`
	Description string   `json:"description"`
	Person      *Person  `json:"person"`
	Color       *Color   `json:"colors"`
	Links       []*Link  `json:"links"`
	Legal       *Legal   `json:"legal"`
	CustomPages []string `json:"custom_pages"`
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
	Link        string `json:"link"`
}

type Color struct {
	Background *BackgroundColor `json:"background"`
	Button     *ButtonColor     `json:"buttons"`
	Text       string           `json:"text"`
	TagHover   string           `json:"tag_hover"`
}

type BackgroundColor struct {
	Type   string `json:"type"`
	Angle  uint   `json:"angle"`
	Colors []struct {
		Color    string `json:"color"`
		Position uint   `json:"position"`
	} `json:"colors"`
}

type ButtonColor struct {
	Text            string `json:"text"`
	TextHover       string `json:"text_hover"`
	Background      string `json:"background"`
	BackgroundHover string `json:"background_hover"`
}

type Link struct {
	Link    string `json:"link"`
	Content string `json:"content"`
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
	return c.Color.GetTextColor()
}

type CustomPage struct {
	Title       string           `json:"title"`
	URI         string           `json:"uri"`
	Image       string           `json:"image"`
	Description string           `json:"description"`
	Color       *Color           `json:"colors"`
	Content     []*CustomContent `json:"content"`
}

type CustomContent struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Content interface {
	Get() template.HTML
}

func (c *Config) LoadCustomPages() ([]*CustomPage, error) {
	if c.CustomPages == nil {
		return nil, nil
	}
	var pages []*CustomPage
	for _, cp := range c.CustomPages {
		b, err := os.ReadFile(cp)
		if err != nil {
			return nil, err
		}
		var p *CustomPage
		err = json.Unmarshal(b, p)
		if err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	return pages, nil
}

func (t *Color) GetTextColor() template.CSS {
	return template.CSS("--text-color: " + t.Text + ";")
}

func (b *ButtonColor) GetTextColor() template.CSS {
	return template.CSS("--text-color: " + b.Text + ";--text-color-hover: " + b.TextHover + ";")
}

func (b *ButtonColor) GetBackground() template.CSS {
	return template.CSS("--background: " + b.Background + ";--background-hover: " + b.BackgroundHover + ";")
}

func (t *Color) GetTagColor() template.CSS {
	return template.CSS("--tag-hover: " + t.TagHover + ";")
}

func (c *CustomContent) Get() template.HTML {
	if c.Type == TitleContentType {
		return template.HTML("<h2>" + c.Content + "</h2>")
	} else if c.Type == SubtitleContentType {
		return template.HTML("<h3>" + c.Content + "</h3>")
	} else if c.Type == ParagraphContentType {
		return template.HTML("<p>" + c.Content + "</p>")
	} else if c.Type == ListContentType {
		v := ""
		for _, s := range strings.Split(c.Content, "--") {
			v += "<li>" + strings.Trim(s, " ") + "</li>"
		}
		return template.HTML("<ul>" + v + "</ul>")
	} else if c.Type == OrderedListContentType {
		v := ""
		for _, s := range strings.Split(c.Content, "--") {
			v += "<li>" + strings.Trim(s, " ") + "</li>"
		}
		return template.HTML("<ol>" + v + "</ol>")
	}
	slog.Warn("Unknown type", "type", c.Type, "value", c.Content)
	return ""
}
