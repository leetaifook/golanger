package web

import (
	"fmt"
	"golanger/session"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

type Page struct {
	*Site
	DefaultController   interface{}
	NotFoundtController interface{}
	Controller          map[string]interface{}
	CurrentController   string
	CurrentAction       string
	Template            string
	TemplateFunc        template.FuncMap
	*Config
	*Document
	pageRLock sync.RWMutex
	pageLock  sync.Mutex
	globalTpl *template.Template
}

type PageParam struct {
	MaxFormSize   int64
	CookieName    string
	Expires       int
	TimerDuration time.Duration
}

func NewPage(param PageParam) *Page {
	if param.MaxFormSize <= 0 {
		param.MaxFormSize = 2 << 20 // 2MB => 2的20次方 乘以 2 =》 2 * 1024 * 1024
	}

	return &Page{
		Site: &Site{
			Base: &Base{
				MAX_FORM_SIZE: param.MaxFormSize,
				Session:       session.New(param.CookieName, param.Expires, param.TimerDuration),
			},
			Version: strconv.Itoa(time.Now().Year()),
		},
		Controller:   map[string]interface{}{},
		TemplateFunc: template.FuncMap{},
		Config:       NewConfig(),
		Document: &Document{
			Css:  map[string]string{},
			Js:   map[string]string{},
			Img:  map[string]string{},
			Func: template.FuncMap{},
		},
	}
}

func (p *Page) Init() {
	p.Site.Init()
	p.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (p *Page) Reset() {
	p.reset(false)
}

func (p *Page) reset(update bool) {
	if update {
		if p.Site.Base.SupportSession != p.Config.SupportSession {
			p.Site.Base.SupportSession = p.Config.SupportSession
		}

		if p.Document.Theme != p.Config.Theme {
			p.Document.Theme = p.Config.Theme
		}

		if p.Site.Root == p.Config.SiteRoot {
			return
		} else {
			p.SetDefaultController(p.GetController(p.Config.IndexDirectory))
			p.UpdateController(p.Site.Root, p.Config.SiteRoot, p.DefaultController)
			p.Site.Root = p.Config.SiteRoot
			p.Document.Static = p.Site.Root + p.Config.StaticDirectory[2:]
		}
	} else {
		p.Site.Base.SupportSession = p.Config.SupportSession
		p.Document.Theme = p.Config.Theme
		p.Site.Root = p.Config.SiteRoot
		p.Document.Static = p.Site.Root + p.Config.StaticDirectory[2:]
		p.SetDefaultController(p.GetController(p.Config.IndexDirectory))
		p.RegisterController(p.Site.Root, p.DefaultController)
		p.globalTpl = template.New("globalTpl").Funcs(p.TemplateFunc)
	}

	siteRootRightTrim := p.Site.Root[:len(p.Site.Root)-1]

	if globalCssFi, err := os.Stat(p.Config.StaticCssDirectory + "/global/"); err == nil && globalCssFi.IsDir() {
		DcssPath := p.Config.StaticCssDirectory + "global/"
		p.Document.Css["global"] = siteRootRightTrim + DcssPath[1:]
		if _, err := os.Stat(DcssPath + "global.css"); err == nil {
			p.Document.GlobalCssFile = p.Document.Css["global"] + "global.css"
		}
	}

	if globalJsFi, err := os.Stat(p.Config.StaticJsDirectory + "/global/"); err == nil && globalJsFi.IsDir() {
		DjsPath := p.Config.StaticJsDirectory + "global/"
		p.Document.Js["global"] = siteRootRightTrim + DjsPath[1:]
		if _, err := os.Stat(DjsPath + "global.js"); err == nil {
			p.Document.GlobalJsFile = p.Document.Js["global"] + "global.js"
		}
	}

	if globalImgFi, err := os.Stat(p.Config.StaticImgDirectory + "/global/"); err == nil && globalImgFi.IsDir() {
		DimgPath := p.Config.StaticImgDirectory + "global/"
		p.Document.Img["global"] = siteRootRightTrim + DimgPath[1:]
	}

	if t, _ := p.globalTpl.ParseGlob(p.Config.TemplateDirectory + p.Config.ThemeDirectory + p.Config.TemplateGlobalDirectory + p.Config.TemplateGlobalFile); t != nil {
		p.globalTpl = t
	}
}

func (p *Page) SetDefaultController(c interface{}) *Page {
	p.DefaultController = c

	return p
}

func (p *Page) SetNotFoundController(c interface{}) *Page {
	p.NotFoundtController = c

	return p
}

func (p *Page) SetTemplate(template string) *Page {
	p.Template = template

	return p
}

func (p *Page) RegisterController(relUrlPath string, i interface{}) *Page {
	if _, ok := p.Controller[relUrlPath]; !ok {
		p.Controller[relUrlPath] = i
	}

	return p
}

func (p *Page) UpdateController(oldUrlPath, relUrlPath string, i interface{}) *Page {
	delete(p.Controller, oldUrlPath)
	p.Controller[relUrlPath] = i

	return p
}

func (p *Page) GetController(urlPath string) (i interface{}) {
	var relUrlPath string
	if strings.HasPrefix(urlPath, p.Site.Root) {
		relUrlPath = urlPath[len(p.Site.Root):]
	} else {
		relUrlPath = urlPath
	}

	i, ok := p.Controller[relUrlPath]
	if !ok {
		i = p.NotFoundtController
	}

	return
}

func (p *Page) AddTemplateFunc(name string, i interface{}) {
	_, ok := p.TemplateFunc[name]
	if !ok {
		p.TemplateFunc[name] = i
	} else {
		fmt.Println("func:" + name + " be added,do not reepeat to add")
	}
}

func (p *Page) DelTemplateFunc(name string) {
	if _, ok := p.TemplateFunc[name]; ok {
		delete(p.TemplateFunc, name)
	}
}

func (p *Page) HandleFavicon() {
	http.HandleFunc(p.Site.Root+"favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		staticPath := p.Config.StaticDirectory + p.Config.ThemeDirectory + "favicon.ico"
		http.ServeFile(w, r, staticPath)
	})
}

func (p *Page) HandleStatic() {
	http.HandleFunc(p.Document.Static, func(w http.ResponseWriter, r *http.Request) {
		staticPath := p.Config.StaticDirectory + r.URL.Path[len(p.Document.Static):]
		http.ServeFile(w, r, staticPath)
	})
}

func (p *Page) handleRoute() {
	http.HandleFunc(p.Site.Root, func(w http.ResponseWriter, r *http.Request) {
		p.pageLock.Lock()
		p.Site.Base.Request = r
		p.Site.Base.ResponseWriter = w

		if p.Config.Reload() {
			p.reset(true)
		}

		p.Site.Base.Cookie = r.Cookies()

		urlPath, fileName := filepath.Split(r.URL.Path)
		if urlPath == p.Site.Root {
			urlPath = p.Site.Root + p.Config.IndexDirectory
		}

		if fileName == "" {
			fileName = p.Config.IndexPage
		}

		p.Template = urlPath[len(p.Site.Root):] + fileName
		fileExt := filepath.Ext(fileName)
		fileNameNoExt := fileName[:len(fileName)-len(fileExt)]
		methodName := strings.Replace(strings.Title(strings.Replace(fileNameNoExt, "_", " ", -1)), " ", "", -1)
		siteRootRightTrim := p.Site.Root[:len(p.Site.Root)-1]

		if cssFi, err := os.Stat(p.Config.StaticCssDirectory + urlPath); err == nil && cssFi.IsDir() {
			cssPath := strings.Trim(urlPath, "/")
			DcssPath := p.Config.StaticCssDirectory + cssPath + "/"
			p.Document.Css[cssPath] = siteRootRightTrim + DcssPath[1:]
			if _, err := os.Stat(DcssPath + "global.css"); err == nil {
				p.Document.GlobalIndexCssFile = p.Document.Css[cssPath] + "global.css"
			}

			if _, err := os.Stat(DcssPath + fileNameNoExt + ".css"); err == nil {
				p.Document.IndexCssFile = p.Document.Css[cssPath] + fileNameNoExt + ".css"
			}

		}

		if jsFi, err := os.Stat(p.Config.StaticJsDirectory + urlPath); err == nil && jsFi.IsDir() {
			jsPath := strings.Trim(urlPath, "/")
			DjsPath := p.Config.StaticJsDirectory + jsPath + "/"
			p.Document.Js[jsPath] = siteRootRightTrim + DjsPath[1:]
			if _, err := os.Stat(DjsPath + "global.js"); err == nil {
				p.Document.GlobalIndexJsFile = p.Document.Js[jsPath] + "global.js"
			}

			if _, err := os.Stat(DjsPath + fileNameNoExt + ".js"); err == nil {
				p.Document.IndexJsFile = p.Document.Js[jsPath] + fileNameNoExt + ".js"
			}
		}

		if imgFi, err := os.Stat(p.Config.StaticImgDirectory + urlPath); err == nil && imgFi.IsDir() {
			imgPath := strings.Trim(urlPath, "/")
			DimgPath := p.Config.StaticImgDirectory + imgPath + "/"
			p.Document.Img[imgPath] = siteRootRightTrim + DimgPath[1:]
		}

		if p.Config.AutoGenerateHtml {
			p.Document.GenerateHtml = true
		}

		p.CurrentController = urlPath[len(p.Site.Root):]
		p.CurrentAction = methodName

		pageController := p.GetController(p.CurrentController)
		rv := reflect.ValueOf(pageController)
		rt := rv.Type()
		if _, ok := rt.MethodByName("Init"); ok {
			rv.MethodByName("Init").Call([]reflect.Value{})
		}

		if _, ok := rt.MethodByName(p.CurrentAction); ok && p.CurrentAction != "Init" && p.Document.Close == false {
			rv.MethodByName(p.CurrentAction).Call([]reflect.Value{})
		} else {
			if !strings.Contains(rt.String(), "Page404") {
				notFountRV := reflect.ValueOf(p.NotFoundtController)
				notFountRV.MethodByName("Init").Call([]reflect.Value{})
			}
		}

		p.pageLock.Unlock()
		p.pageRLock.RLock()

		if p.Document.Close == false && p.Document.Hide == false {
			if tplFi, err := os.Stat(p.Config.TemplateDirectory + p.Config.ThemeDirectory + p.Template); err == nil {
				globalTemplate, _ := p.globalTpl.Clone()
				if pageTemplate, err := globalTemplate.New(filepath.Base(p.Template)).ParseFiles(p.Config.TemplateDirectory + p.Config.ThemeDirectory + p.Template); err == nil {
					templateVar := map[string]interface{}{
						"G":        p.Base.GET,
						"P":        p.Base.POST,
						"C":        p.Base.COOKIE,
						"S":        p.Base.SESSION,
						"Siteroot": p.Site.Root,
						"Version":  p.Site.Version,
						"Template": p.Template,
						"D":        p.Document,
						"Config":   p.Config.M,
					}

					if p.Document.GenerateHtml {
						htmlFile := p.Config.StaticDirectory + p.Config.HtmlDirectory + urlPath + fileName
						htmlDir := filepath.Dir(htmlFile)
						if htmlDirFi, err := os.Stat(htmlDir); err != nil || !htmlDirFi.IsDir() {
							os.MkdirAll(htmlDir, 0777)
						}

						var doWrite bool
						if p.Config.AutoGenerateHtml {
							if p.Config.AutoGenerateHtmlCycleTime <= 0 {
								doWrite = true
							} else {
								if htmlFi, err := os.Stat(htmlFile); err != nil {
									doWrite = true
								} else {
									switch {
									case tplFi.ModTime().Unix() >= htmlFi.ModTime().Unix():
										doWrite = true
									case tplFi.ModTime().Unix() >= htmlFi.ModTime().Unix():
										doWrite = true
									case time.Now().Unix()-htmlFi.ModTime().Unix() >= p.Config.AutoGenerateHtmlCycleTime:
										doWrite = true
									default:
										globalTplFi, err := os.Stat(p.Config.TemplateDirectory + p.Config.ThemeDirectory + p.Config.TemplateGlobalDirectory)
										if err == nil {
											if globalTplFi.ModTime().Unix() >= htmlFi.ModTime().Unix() {
												doWrite = true
											}
										}
									}
								}
							}
						}

						if doWrite {
							if file, err := os.OpenFile(htmlFile, os.O_CREATE|os.O_WRONLY, 0777); err == nil {
								templateVar["Siteroot"] = p.Config.SiteRoot + htmlDir + "/"
								p.pageLock.Lock()
								pageTemplate.Execute(file, templateVar)
								p.pageLock.Unlock()
							}
						}

						if p.Config.AutoJumpToHtml {
							p.pageLock.Lock()
							http.Redirect(w, r, p.Site.Root+htmlFile[2:], http.StatusFound)
							p.pageLock.Unlock()
						} else {
							p.pageLock.Lock()
							err := pageTemplate.Execute(w, templateVar)
							p.pageLock.Unlock()
							if err != nil {
								log.Println(err)
							}
						}
					} else {
						p.pageLock.Lock()
						err := pageTemplate.Execute(w, templateVar)
						if err != nil {
							log.Println(err)
							w.Write([]byte(fmt.Sprint(err)))
						}

						p.pageLock.Unlock()
					}
				} else {
					log.Println(err)
					p.pageLock.Lock()
					w.Write([]byte(fmt.Sprint(err)))
					p.pageLock.Unlock()
				}
			}
		}
		p.pageRLock.RUnlock()

		p.pageLock.Lock()
		p.Document.Reset()
		p.pageLock.Unlock()
	})
}

func (p *Page) ListenAndServe(addr string) {
	p.handleRoute()

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
