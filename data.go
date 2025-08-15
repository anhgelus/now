package main

import (
	"encoding/json"
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/anhgelus/golatt"
	"html/template"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	regexExternalLink = regexp.MustCompile(`https?://`)
)

type ConfigData interface {
	GetTextColor() template.CSS
	GetBackground() template.CSS
	GetBackgroundImage() template.CSS
	IsCustomPage() bool
}

type Config struct {
	Image       string   `json:"image" toml:"image"`
	Description string   `json:"description" toml:"description"`
	Person      *Person  `json:"person" toml:"person"`
	Color       *Color   `json:"colors" toml:"colors"`
	Links       []*Link  `json:"links" toml:"links"`
	Legal       string   `json:"legal" toml:"legal"`
	RelMeLinks  []string `json:"rel_me_links" toml:"rel_me_links"`
	CustomPages []string `json:"custom_pages" toml:"custom_pages"`
	folder      string
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

func getImage(s string) string {
	if regexExternalLink.MatchString(s) {
		return s
	}
	return golatt.GetStaticPath(s)
}

func (c *Config) GetBackground() template.CSS {
	return c.Color.GetBackground()
}

func (c *Config) GetBackgroundImage() template.CSS {
	return template.CSS("--background-image: url(" + getImage(c.Image) + ");")
}

func (c *Config) GetTextColor() template.CSS {
	return c.Color.GetTextColor()
}

func (c *Config) IsCustomPage() bool {
	return false
}

var legalContent template.HTML

func (c *Config) GetLegal() (template.HTML, error) {
	if legalContent == "" {
		b, err := os.ReadFile(c.folder + c.Legal)
		if err != nil {
			return "", err
		}
		legalContent = template.HTML(b)
	}
	return legalContent, nil
}

type CustomPage struct {
	Title       string `json:"title" toml:"title"`
	URI         string `json:"uri" toml:"uri"`
	Image       string `json:"image" toml:"image"`
	Description string `json:"description" toml:"description"`
	Color       *Color `json:"colors" toml:"colors"`
	Content     string `json:"content" toml:"content"`
	folder      string
}

func (c *Config) LoadCustomPages() ([]*CustomPage, error) {
	if c.CustomPages == nil {
		println("null")
		return nil, nil
	}
	var pages []*CustomPage
	for _, cp := range c.CustomPages {
		b, err := os.ReadFile(c.folder + cp)
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

		p.folder = getFolder(c.folder + cp)

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
	return template.CSS("--background-image: url(" + getImage(p.Image) + ");")
}

func (p *CustomPage) GetBackground() template.CSS {
	return p.Color.GetBackground()
}

func (p *CustomPage) IsCustomPage() bool {
	return true
}

var contentsMap = map[string]template.HTML{}

func (p *CustomPage) GetContent() (template.HTML, error) {
	res, ok := contentsMap[p.URI]
	if !ok {
		b, err := os.ReadFile(p.folder + p.Content)
		if err != nil {
			return "", err
		}
		res = template.HTML(b)
		contentsMap[p.URI] = res
	}
	return res, nil
}
