package web

import (
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
)

type Config struct {
	SupportSession            bool                   "SupportSession"
	AutoGenerateHtml          bool                   "AutoGenerateHtml"
	AutoGenerateHtmlCycleTime int64                  "AutoGenerateHtmlCycleTime"
	AutoJumpToHtml            bool                   "AutoJumpToHtml"
	Debug                     bool                   "Debug"
	StaticDirectory           string                 "StaticDirectory"
	ThemeDirectory            string                 "ThemeDirectory"
	Theme                     string                 "Theme"
	StaticCssDirectory        string                 "StaticCssDirectory"
	StaticJsDirectory         string                 "StaticJsDirectory"
	StaticImgDirectory        string                 "StaticImgDirectory"
	HtmlDirectory             string                 "HtmlDirectory"
	TemplateDirectory         string                 "TemplateDirectory"
	TemplateGlobalDirectory   string                 "TemplateGlobalDirectory"
	TemplateGlobalFile        string                 "TemplateGlobalFile"
	TemporaryDirectory        string                 "TemporaryDirectory"
	UploadDirectory           string                 "UploadDirectory"
	IndexDirectory            string                 "IndexDirectory"
	IndexPage                 string                 "IndexPage"
	SiteRoot                  string                 "SiteRoot"
	Environment               map[string]string      "Environment"
	Database                  map[string]string      "Database"
	M                         map[string]interface{} "Custom"
	configPath                string
	configLastModTime         int64
}

func NewConfig() Config {
	return Config{
		TemplateDirectory:       "./view/",
		TemporaryDirectory:      "./tmp/",
		StaticDirectory:         "./static/",
		ThemeDirectory:          "theme/",
		Theme:                   "default",
		StaticCssDirectory:      "css/",
		StaticJsDirectory:       "js/",
		StaticImgDirectory:      "img/",
		HtmlDirectory:           "html/",
		UploadDirectory:         "upload/",
		TemplateGlobalDirectory: "_global/",
		TemplateGlobalFile:      "*",
		IndexDirectory:          "index/",
		IndexPage:               "index.html",
		SiteRoot:                "/",
		Environment:             map[string]string{},
		Database:                map[string]string{},
	}
}

func (c *Config) Load(configPath string) {
	yamlData, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	err = goyaml.Unmarshal(yamlData, c)
	if err != nil {
		panic(err)
	}

	c.UploadDirectory = c.StaticDirectory + c.UploadDirectory
	c.ThemeDirectory = c.ThemeDirectory + c.Theme + "/"
	c.StaticCssDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticCssDirectory
	c.StaticJsDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticJsDirectory
	c.StaticImgDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticImgDirectory

	c.configPath = configPath
	yamlFi, _ := os.Stat(configPath)
	c.configLastModTime = yamlFi.ModTime().Unix()
}

func (c *Config) Reload() bool {
	var b bool
	configPath := c.configPath
	yamlFi, _ := os.Stat(configPath)
	if yamlFi.ModTime().Unix() > c.configLastModTime {
		yamlData, _ := ioutil.ReadFile(configPath)
		*c = NewConfig()
		goyaml.Unmarshal(yamlData, c)
		c.configPath = configPath
		c.configLastModTime = yamlFi.ModTime().Unix()
		c.UploadDirectory = c.StaticDirectory + c.UploadDirectory
		c.ThemeDirectory = c.ThemeDirectory + c.Theme + "/"
		c.StaticCssDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticCssDirectory
		c.StaticJsDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticJsDirectory
		c.StaticImgDirectory = c.StaticDirectory + c.ThemeDirectory + c.StaticImgDirectory
		b = true
	}

	return b
}
