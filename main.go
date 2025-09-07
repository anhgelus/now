package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/anhgelus/golatt"
)

var (
	//go:embed templates
	templates embed.FS
	//go:embed dist
	assets embed.FS
)

var (
	domain        string
	configPath    string
	publicDirPath string
	dev           = false
	generateToml  bool
	generateJson  bool
	port          = 80
)

func init() {
	flag.StringVar(&domain, "domain", "", "domain to use")
	flag.StringVar(&configPath, "config", "", "config to use")
	flag.StringVar(&publicDirPath, "public-dir", "", "public directory to use, default is 'public' inside the folder of your config")
	flag.BoolVar(&dev, "dev", dev, "dev mode enabled")
	flag.BoolVar(&generateJson, "generate-json-config", false, "generate a config example")
	flag.BoolVar(&generateToml, "generate-toml-config", false, "generate a config example")
	flag.IntVar(&port, "port", port, "set the port to use")
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

	cfg.folder = getFolder(configPath)

	customPages, err := cfg.LoadCustomPages()
	if err != nil {
		panic(err)
	}

	publicFolder := cfg.folder + "public"
	if len(publicDirPath) != 0 {
		publicFolder = publicDirPath
	}

	var g *golatt.Golatt
	if dev {
		g = golatt.New(
			golatt.UsableEmbedFS("templates", templates),
			os.DirFS(publicFolder),
			os.DirFS("dist"),
		)
	} else {
		g = golatt.New(
			golatt.UsableEmbedFS("templates", templates),
			os.DirFS(publicFolder),
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
		"Legal information about "+cfg.Person.Name+"'s bio",
		&cfg,
	).Handle()
	g.NewTemplate("now",
		"/now",
		"Now",
		"",
		""+cfg.Person.Name+"'s now",
		&cfg,
	).Handle()

	for _, cp := range customPages {
		slog.Info("Creating custom page...", "title", cp.Title, "uri", cp.URI)
		g.NewTemplate("custom_page",
			cp.URI,
			cp.Title,
			cp.Image,
			cp.Description,
			cp,
		).Handle()
	}

	g.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		g.Render(w, "404", &golatt.TemplateData{
			Title: "Not found :(",
			SEO: &golatt.SeoData{
				URL:         r.URL.Path,
				Description: "Not found",
			},
			Data: &cfg,
		})
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	g.TemplateFuncMap = template.FuncMap{
		"getImage": getImage,
		"getRings": func() []*Ring { return cfg.Rings },
		"getFont":  func() string { return cfg.Font },
	}

	host := fmt.Sprintf(":%d", port)
	if dev {
		if port != 80 {
			g.StartServer(host)
		} else {
			g.StartServer(":8000")
		}
	} else {
		g.StartServer(host)
	}
}

func getFolder(path string) string {
	if !strings.Contains(path, "/") {
		return ""
	}
	sp := strings.Split(path, "/")
	folder := strings.Join(sp[1:len(sp)-1], "/")
	if path[0] != '/' && !strings.HasPrefix(path, "./") {
		if len(folder) == 0 {
			folder = sp[0]
		} else {
			folder = sp[0] + "/" + folder
		}
	}
	return folder + "/"
}

func generateConfigFile(isToml bool) {
	cfg := Config{
		Image:       "wallpaper.webp",
		Description: "I am a beautiful description!",
		Person: &Person{
			Name:     "John Doe",
			Pronouns: "any",
			Image:    "pfp.webp",
			Now: []*Now{
				{"Hello", "World", "", ""},
				{"I am", "a tag", "", ""},
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
		Legal:       "legal.html",
		CustomPages: []string{"custom.json"},
		RelMeLinks:  []string{"https://foo.example.org/@bar"},
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
	fmt.Println(string(b))
}
