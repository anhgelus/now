package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
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
	ButtonsContentType     = "links"
)

type ConfigData interface {
	GetTextColor() template.CSS
	GetBackground() template.CSS
	GetBackgroundImage() template.CSS
}

type Config struct {
	Image       string   `json:"image" toml:"image"`
	Description string   `json:"description" toml:"description"`
	Person      *Person  `json:"person" toml:"person"`
	Color       *Color   `json:"colors" toml:"colors"`
	Links       []*Link  `json:"links" toml:"links"`
	Legal       *Legal   `json:"legal" toml:"legal"`
	CustomPages []string `json:"custom_pages" toml:"custom_pages"`
}

type Person struct {
	Name     string `json:"name" toml:"name"`
	Pronouns string `json:"pronouns" toml:"pronouns"`
	Image    string `json:"image" toml:"image"`
	Tags     []*Tag `json:"tags" toml:"tags"`
}

type Tag struct {
	Name        string `json:"name" toml:"name"`
	Description string `json:"description" toml:"description"`
	Link        string `json:"link" toml:"link"`
}

type Color struct {
	Background *BackgroundColor `json:"background" toml:"background"`
	Button     *ButtonColor     `json:"buttons" toml:"buttons"`
	Text       string           `json:"text" toml:"text"`
	TagHover   string           `json:"tag_hover" toml:"tag_hover"`
}

type BackgroundColor struct {
	Type   string `json:"type" toml:"type"`
	Angle  uint   `json:"angle" toml:"angle"`
	Colors []struct {
		Color    string `json:"color" toml:"color"`
		Position uint   `json:"position" toml:"position"`
	} `json:"colors" toml:"colors"`
}

type ButtonColor struct {
	Text            string `json:"text" toml:"text"`
	TextHover       string `json:"text_hover" toml:"text_hover"`
	Background      string `json:"background" toml:"background"`
	BackgroundHover string `json:"background_hover" toml:"background_hover"`
}

type Link struct {
	Link    string `json:"link" toml:"link"`
	Content string `json:"content" toml:"content"`
}

type Legal struct {
	LegalInformationLink string   `json:"legal_information_link" toml:"legal_information_link"`
	ImagesSource         []string `json:"images_source" toml:"images_source"`
	FontSource           string   `json:"font_source" toml:"font_source"`
}

func (c *Config) GetBackground() template.CSS {
	return c.Color.GetBackground()
}

func (c *Config) GetBackgroundImage() template.CSS {
	return template.CSS("--background-image: url(" + golatt.GetStaticPath(c.Image) + ");")
}

func (c *Config) GetTextColor() template.CSS {
	return c.Color.GetTextColor()
}

type CustomPage struct {
	Title       string           `json:"title" toml:"title"`
	URI         string           `json:"uri" toml:"uri"`
	Image       string           `json:"image" toml:"image"`
	Description string           `json:"description" toml:"description"`
	Color       *Color           `json:"colors" toml:"colors"`
	Content     []*CustomContent `json:"content" toml:"content"`
}

type CustomContent struct {
	Type    string `json:"type" toml:"type"`
	Content string `json:"content" toml:"content"`
}

type Content interface {
	Get() template.HTML
}

func (c *Config) LoadCustomPages() ([]*CustomPage, error) {
	if c.CustomPages == nil {
		println("null")
		return nil, nil
	}
	var pages []*CustomPage
	for _, cp := range c.CustomPages {
		b, err := os.ReadFile(cp)
		if err != nil {
			return nil, err
		}
		var p CustomPage
		if strings.HasSuffix(cp, ".json") {
			err = json.Unmarshal(b, &p)
		} else if strings.HasSuffix(cp, ".toml") {
			err = toml.Unmarshal(b, &p)
		} else {
			return nil, errors.New("custom page file must be .json or .toml")
		}
		if err != nil {
			return nil, err
		}
		pages = append(pages, &p)
	}
	return pages, nil
}

func (t *Color) GetTextColor() template.CSS {
	return template.CSS("--text-color: " + t.Text + ";")
}

func (t *Color) GetBackground() template.CSS {
	bg := t.Background
	css := "background: " + bg.Type + "-gradient("
	if bg.Type == "linear" {
		css += strconv.Itoa(int(bg.Angle)) + "deg,"
	}
	for _, c := range bg.Colors {
		css += c.Color + " " + strconv.Itoa(int(c.Position)) + "%,"
	}
	return template.CSS(css[:len(css)-1] + ");")
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

func (p *CustomPage) GetTextColor() template.CSS {
	return p.Color.GetTextColor()
}

func (p *CustomPage) GetBackgroundImage() template.CSS {
	return template.CSS("--background-image: url(" + golatt.GetStaticPath(p.Image) + ");")
}

func (p *CustomPage) GetBackground() template.CSS {
	return p.Color.GetBackground()
}

func (p *CustomPage) GetContent() template.HTML {
	var res template.HTML
	for _, c := range p.Content {
		res += c.Get(p)
	}
	return res
}

func (c *CustomContent) Get(p *CustomPage) template.HTML {
	if c.Type == TitleContentType {
		return template.HTML("<h2>" + c.Content + "</h2>")
	} else if c.Type == SubtitleContentType {
		return template.HTML("<h3>" + c.Content + "</h3>")
	} else if c.Type == ParagraphContentType {
		return template.HTML("<p>" + c.Content + "</p>")
	} else if c.Type == ListContentType {
		v := ""
		for _, s := range strings.Split(c.Content, "--") {
			if len(strings.Trim(s, " ")) == 0 {
				continue
			}
			v += "<li>" + strings.Trim(s, " ") + "</li>"
		}
		return template.HTML("<ul>" + v + "</ul>")
	} else if c.Type == OrderedListContentType {
		v := ""
		for _, s := range strings.Split(c.Content, "--") {
			if len(strings.TrimSpace(s)) == 0 {
				continue
			}
			v += "<li>" + strings.TrimSpace(s) + "</li>"
		}
		return template.HTML("<ol>" + v + "</ol>")
	} else if c.Type == ButtonsContentType {
		// [Bonsoir](/hello) -- [Bonjour](/not_hello)
		v := ""
		for _, s := range strings.Split(c.Content, "--") {
			if len(strings.TrimSpace(s)) == 0 {
				continue
			}
			sp := strings.Split(s, "](")
			if len(sp) != 2 {
				slog.Warn("Invalid button", "s", s)
				continue
			}
			url := strings.TrimSpace(sp[1])
			v += fmt.Sprintf(
				`<div class="link"><a href="%s">%s</a></div>`,
				url[:len(url)-1],
				strings.TrimSpace(sp[0])[1:],
			)
		}
		return template.HTML(fmt.Sprintf(
			`<nav class="links" style="%s">%s</nav>`,
			p.Color.Button.GetBackground()+p.Color.Button.GetTextColor(),
			v,
		))
	}
	slog.Warn("Unknown type", "type", c.Type, "value", c.Content)
	return ""
}
