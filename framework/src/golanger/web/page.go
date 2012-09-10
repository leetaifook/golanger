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
	header              map[string][]string
	Controller          map[string]interface{}
	DefaultController   interface{}
	NotFoundtController interface{}
	CurrentPath         string
	CurrentFileName     string
	CurrentController   string
	CurrentAction       string
	Template            string
	TemplateFunc        template.FuncMap
	*Config
	Document
	supportStatic bool
	rmutex        sync.RWMutex
	mutex         sync.Mutex
	globalTpl     *template.Template
}

type PageParam struct {
	MaxFormSize   int64
	CookieName    string
	Expires       int
	TimerDuration time.Duration
}

func NewPage(param PageParam) Page {
	if param.MaxFormSize <= 0 {
		param.MaxFormSize = 2 << 20 // 2MB => 2的20次方 乘以 2 =》 2 * 1024 * 1024
	}

	return Page{
		Site: &Site{
			Base: &Base{
				MAX_FORM_SIZE: param.MaxFormSize,
				Session:       session.New(param.CookieName, param.Expires, param.TimerDuration),
			},
			Version: strconv.Itoa(time.Now().Year()),
		},
		header:       map[string][]string{},
		Controller:   map[string]interface{}{},
		TemplateFunc: template.FuncMap{},
		Config:       NewConfig(),
		Document: Document{
			Css:  map[string]string{},
			Js:   map[string]string{},
			Img:  map[string]string{},
			Func: template.FuncMap{},
		},
	}
}

func (p *Page) Init(w http.ResponseWriter, r *http.Request) {
	p.Site.Init(w, r)

	if p.header != nil || len(p.header) > 0 {
		for t, s := range p.header {
			for _, v := range s {
				w.Header().Add(t, v)
			}
		}
	}
}

func (p *Page) SetDefaultController(i interface{}) *Page {
	p.DefaultController = i

	return p
}

func (p *Page) SetNotFoundController(i interface{}) *Page {
	p.NotFoundtController = i

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

func (p *Page) GetController(urlPath string) interface{} {
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

	return i
}

func (p *Page) AddHeader(k, v string) {
	if _, ok := p.header[k]; ok {
		p.header[k] = append(p.header[k], v)
	} else {
		p.header[k] = []string{v}
	}
}

func (p *Page) DelHeader(k string) {
	delete(p.header, k)
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

func (p *Page) Load(configPath string) {
	p.Config.Load(configPath)
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

		if p.Document.Static != p.Config.SiteRoot+p.Config.StaticDirectory[2:] {
			p.Document.Static = p.Config.SiteRoot + p.Config.StaticDirectory[2:]
		}

		if p.Site.Root == p.Config.SiteRoot {
			return
		} else {
			p.SetDefaultController(p.GetController(p.Config.IndexDirectory))
			p.UpdateController(p.Site.Root, p.Config.SiteRoot, p.DefaultController)
			p.Site.Root = p.Config.SiteRoot
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

func (p *Page) setCurrentInfo(path string) {
	urlPath, fileName := filepath.Split(path)
	if urlPath == p.Site.Root {
		urlPath = p.Site.Root + p.Config.IndexDirectory
	}

	if fileName == "" {
		fileName = p.Config.IndexPage
	}

	p.CurrentPath = urlPath
	p.CurrentFileName = fileName
	p.CurrentController = urlPath[len(p.Site.Root):]
	p.CurrentAction = strings.Replace(strings.Title(strings.Replace(p.CurrentFileName[:len(p.CurrentFileName)-len(filepath.Ext(p.CurrentFileName))], "_", " ", -1)), " ", "", -1)
}

func (p *Page) routeController(i interface{}, w http.ResponseWriter, r *http.Request) {
	p.mutex.Lock()
	p.Site.Base.Cookie = r.Cookies()
	p.setCurrentInfo(r.URL.Path)
	p.Template = p.CurrentController + p.CurrentFileName
	p.mutex.Unlock()

	p.rmutex.RLock()
	pageOriController := p.GetController(p.CurrentPath)
	rv := reflect.ValueOf(pageOriController)
	p.rmutex.RUnlock()

	rvw, rvr := reflect.ValueOf(w), reflect.ValueOf(r)
	rt := rv.Type()
	vpc := reflect.New(rt)
	iv := reflect.ValueOf(i).Elem()
	vpc.Elem().FieldByName("Application").Set(iv)
	tpc := vpc.Type()
	if _, found := tpc.Elem().FieldByName("RW"); found {
		rvarw := vpc.Elem().FieldByName("RW")
		rvarw.Set(rvw)
	}

	if _, found := tpc.Elem().FieldByName("R"); found {
		rvar := vpc.Elem().FieldByName("R")
		rvar.Set(rvr)
	}

	vppc := vpc.Elem().FieldByName("Page")
	ppc := vppc.Interface().(Page)

	if _, ok := tpc.MethodByName(ppc.CurrentAction); ok && ppc.CurrentAction != "Init" {
		if rm, ok := tpc.MethodByName("Init"); ok {
			mt := rm.Type
			switch mt.NumIn() {
			case 2:
				if mt.In(1) == rvr.Type() {
					p.rmutex.RLock()
					vpc.MethodByName("Init").Call([]reflect.Value{rvr})
					p.rmutex.RUnlock()
				} else {
					p.rmutex.RLock()
					vpc.MethodByName("Init").Call([]reflect.Value{rvw})
					p.rmutex.RUnlock()
				}
			case 3:
				p.rmutex.RLock()
				vpc.MethodByName("Init").Call([]reflect.Value{rvw, rvr})
				p.rmutex.RUnlock()
			default:
				p.rmutex.RLock()
				vpc.MethodByName("Init").Call([]reflect.Value{})
				p.rmutex.RUnlock()
			}
		}

		if ppc.Document.Close == false {
			rm, _ := tpc.MethodByName(ppc.CurrentAction)
			mt := rm.Type
			switch mt.NumIn() {
			case 2:
				if mt.In(1) == rvr.Type() {
					p.rmutex.RLock()
					vpc.MethodByName(ppc.CurrentAction).Call([]reflect.Value{rvr})
					p.rmutex.RUnlock()
				} else {
					p.rmutex.RLock()
					vpc.MethodByName(ppc.CurrentAction).Call([]reflect.Value{rvw})
					p.rmutex.RUnlock()
				}
			case 3:
				p.rmutex.RLock()
				vpc.MethodByName(ppc.CurrentAction).Call([]reflect.Value{rvw, rvr})
				p.rmutex.RUnlock()
			default:
				p.rmutex.RLock()
				vpc.MethodByName(ppc.CurrentAction).Call([]reflect.Value{})
				p.rmutex.RUnlock()
			}
		}

		ppc = vppc.Interface().(Page)
	} else {
		if !strings.Contains(tpc.String(), "Page404") {
			notFountRV := reflect.ValueOf(ppc.NotFoundtController)
			notFountRT := notFountRV.Type()
			vnpc := reflect.New(notFountRT)
			vnpc.Elem().FieldByName("Application").Set(iv)
			tnpc := vnpc.Type()
			if _, found := tnpc.Elem().FieldByName("RW"); found {
				rvarw := vnpc.Elem().FieldByName("RW")
				rvarw.Set(rvw)
			}

			if _, found := tnpc.Elem().FieldByName("R"); found {
				rvar := vnpc.Elem().FieldByName("R")
				rvar.Set(rvr)
			}

			vppc = vnpc.Elem().FieldByName("Page")

			if rm, ok := tnpc.MethodByName("Init"); ok {
				mt := rm.Type
				switch mt.NumIn() {
				case 2:
					if mt.In(1) == rvr.Type() {
						p.rmutex.RLock()
						vnpc.MethodByName("Init").Call([]reflect.Value{rvr})
						p.rmutex.RUnlock()
					} else {
						p.rmutex.RLock()
						vnpc.MethodByName("Init").Call([]reflect.Value{rvw})
						p.rmutex.RUnlock()
					}
				case 3:
					p.rmutex.RLock()
					vnpc.MethodByName("Init").Call([]reflect.Value{rvw, rvr})
					p.rmutex.RUnlock()
				default:
					p.rmutex.RLock()
					vnpc.MethodByName("Init").Call([]reflect.Value{})
					p.rmutex.RUnlock()
				}
			}

			ppc = vppc.Interface().(Page)
		}
	}

	vppc.Set(reflect.ValueOf(ppc))

	p.rmutex.RLock()
	if ppc.supportStatic {
		ppc.setStaticDocument()
		ppc.routeTemplate(w, r)
	}
	p.rmutex.RUnlock()

}

func (p *Page) setStaticDocument() {
	fileNameNoExt := p.CurrentFileName[:len(p.CurrentFileName)-len(filepath.Ext(p.CurrentFileName))]
	siteRootRightTrim := p.Site.Root[:len(p.Site.Root)-1]

	if cssFi, err := os.Stat(p.Config.StaticCssDirectory + p.CurrentPath); err == nil && cssFi.IsDir() {
		cssPath := strings.Trim(p.CurrentPath, "/")
		DcssPath := p.Config.StaticCssDirectory + cssPath + "/"
		p.Document.Css[cssPath] = siteRootRightTrim + DcssPath[1:]
		if _, err := os.Stat(DcssPath + "global.css"); err == nil {
			p.Document.GlobalIndexCssFile = p.Document.Css[cssPath] + "global.css"
		}

		if _, err := os.Stat(DcssPath + fileNameNoExt + ".css"); err == nil {
			p.Document.IndexCssFile = p.Document.Css[cssPath] + fileNameNoExt + ".css"
		}

	}

	if jsFi, err := os.Stat(p.Config.StaticJsDirectory + p.CurrentPath); err == nil && jsFi.IsDir() {
		jsPath := strings.Trim(p.CurrentPath, "/")
		DjsPath := p.Config.StaticJsDirectory + jsPath + "/"
		p.Document.Js[jsPath] = siteRootRightTrim + DjsPath[1:]
		if _, err := os.Stat(DjsPath + "global.js"); err == nil {
			p.Document.GlobalIndexJsFile = p.Document.Js[jsPath] + "global.js"
		}

		if _, err := os.Stat(DjsPath + fileNameNoExt + ".js"); err == nil {
			p.Document.IndexJsFile = p.Document.Js[jsPath] + fileNameNoExt + ".js"
		}
	}

	if imgFi, err := os.Stat(p.Config.StaticImgDirectory + p.CurrentPath); err == nil && imgFi.IsDir() {
		imgPath := strings.Trim(p.CurrentPath, "/")
		DimgPath := p.Config.StaticImgDirectory + imgPath + "/"
		p.Document.Img[imgPath] = siteRootRightTrim + DimgPath[1:]
	}
}

func (p *Page) routeTemplate(w http.ResponseWriter, r *http.Request) {
	if p.Config.AutoGenerateHtml {
		p.Document.GenerateHtml = true
	}

	if p.Document.Close == false && p.Document.Hide == false {
		if tplFi, err := os.Stat(p.Config.TemplateDirectory + p.Config.ThemeDirectory + p.Template); err == nil {
			globalTemplate, _ := p.globalTpl.Clone()
			if pageTemplate, err := globalTemplate.New(filepath.Base(p.Template)).ParseFiles(p.Config.TemplateDirectory + p.Config.ThemeDirectory + p.Template); err == nil {
				p.rmutex.RLock()
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
				p.rmutex.RUnlock()

				if p.Document.GenerateHtml {

					htmlFile := p.Config.StaticDirectory + p.Config.HtmlDirectory + p.Site.Root + p.Template
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
							pageTemplate.Execute(file, templateVar)
						}
					}

					if p.Config.AutoJumpToHtml {
						http.Redirect(w, r, p.Site.Root+htmlFile[2:], http.StatusFound)
					} else {
						err := pageTemplate.Execute(w, templateVar)
						if err != nil {
							log.Println(err)
						}
					}
				} else {
					err := pageTemplate.Execute(w, templateVar)
					if err != nil {
						log.Println(err)
						w.Write([]byte(fmt.Sprint(err)))
					}
				}
			} else {
				log.Println(err)
				w.Write([]byte(fmt.Sprint(err)))
			}
		}
	}
}

func (p *Page) HandleFavicon() {
	p.supportStatic = true
	http.HandleFunc(p.Site.Root+"favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		staticPath := p.Config.StaticDirectory + p.Config.ThemeDirectory + "favicon.ico"
		http.ServeFile(w, r, staticPath)
	})
}

func (p *Page) HandleStatic() {
	p.supportStatic = true
	http.HandleFunc(p.Document.Static, func(w http.ResponseWriter, r *http.Request) {
		staticPath := p.Config.StaticDirectory + r.URL.Path[len(p.Document.Static):]
		http.ServeFile(w, r, staticPath)
	})
}

func (p *Page) handleRoute(i interface{}) {
	http.HandleFunc(p.Site.Root, func(w http.ResponseWriter, r *http.Request) {
		p.mutex.Lock()
		if p.Config.Reload() {
			p.reset(true)
		}
		p.mutex.Unlock()

		p.routeController(i, w, r)
	})
}

func (p *Page) ListenAndServe(addr string, i interface{}) {
	p.handleRoute(i)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
