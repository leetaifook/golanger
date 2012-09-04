package main

import (
	"controllers"
	"flag"
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	_ "templateFunc"
	"text/template"
	"time"
)

var (
	addr       = flag.String("addr", ":80", "Server port")
	configPath = flag.String("config", "./config/site.yaml", "site filepath of config")
)

var (
	yamlLastModTime int64
)

func init() {
	flag.Parse()
	os.Chdir(filepath.Dir(os.Args[0]))
	fmt.Println("Listen server address: " + *addr)

	yamlData, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	err = goyaml.Unmarshal(yamlData, &controllers.Page.Config)
	if err != nil {
		panic(err)
	}

	fmt.Println("Read configuration file success, fithpath: " + filepath.Join(filepath.Dir(os.Args[0]), *configPath))
	yamlFi, _ := os.Stat(*configPath)
	yamlLastModTime = yamlFi.ModTime().Unix()

	controllers.Page.Site.Base.SupportSession = controllers.Page.Config.SupportSession
	controllers.Page.Site.Root = controllers.Page.Config.SiteRoot
	controllers.Page.Config.UploadDirectory = controllers.Page.Config.StaticDirectory + controllers.Page.Config.UploadDirectory
	controllers.Page.Config.ThemeDirectory = controllers.Page.Config.ThemeDirectory + controllers.Page.Config.Theme + "/"
	controllers.Page.Config.StaticCssDirectory = controllers.Page.Config.StaticDirectory + controllers.Page.Config.ThemeDirectory + controllers.Page.Config.StaticCssDirectory
	controllers.Page.Config.StaticJsDirectory = controllers.Page.Config.StaticDirectory + controllers.Page.Config.ThemeDirectory + controllers.Page.Config.StaticJsDirectory
	controllers.Page.Config.StaticImgDirectory = controllers.Page.Config.StaticDirectory + controllers.Page.Config.ThemeDirectory + controllers.Page.Config.StaticImgDirectory
	controllers.Page.Document.Static = controllers.Page.Site.Root + controllers.Page.Config.StaticDirectory[2:]
	controllers.Page.Document.Theme = controllers.Page.Config.Theme

	siteRootRightTrim := controllers.Page.Site.Root[:len(controllers.Page.Site.Root)-1]

	if globalCssFi, err := os.Stat(controllers.Page.Config.StaticCssDirectory + "/global/"); err == nil && globalCssFi.IsDir() {
		DcssPath := controllers.Page.Config.StaticCssDirectory + "global/"
		controllers.Page.Document.Css["global"] = siteRootRightTrim + DcssPath[1:]
		if _, err := os.Stat(DcssPath + "global.css"); err == nil {
			controllers.Page.Document.GlobalCssFile = controllers.Page.Document.Css["global"] + "global.css"
		}
	}

	if globalJsFi, err := os.Stat(controllers.Page.Config.StaticJsDirectory + "/global/"); err == nil && globalJsFi.IsDir() {
		DjsPath := controllers.Page.Config.StaticJsDirectory + "global/"
		controllers.Page.Document.Js["global"] = siteRootRightTrim + DjsPath[1:]
		if _, err := os.Stat(DjsPath + "global.js"); err == nil {
			controllers.Page.Document.GlobalJsFile = controllers.Page.Document.Js["global"] + "global.js"
		}
	}

	if globalImgFi, err := os.Stat(controllers.Page.Config.StaticImgDirectory + "/global/"); err == nil && globalImgFi.IsDir() {
		DimgPath := controllers.Page.Config.StaticImgDirectory + "global/"
		controllers.Page.Document.Img["global"] = siteRootRightTrim + DimgPath[1:]
	}

	controllers.Page.SetNotFoundController(&controllers.Page404{controllers.Page})
	controllers.Page.SetDefaultController(controllers.Page.GetController(controllers.Page.Config.IndexDirectory))
	controllers.Page.RegisterController(controllers.Page.Site.Root, controllers.Page.DefaultController)
}

func startApp() {
	http.HandleFunc(controllers.Page.Site.Root+"favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		staticPath := controllers.Page.Config.StaticDirectory + controllers.Page.Config.ThemeDirectory + "favicon.ico"
		http.ServeFile(w, r, staticPath)
	})

	http.HandleFunc(controllers.Page.Document.Static, func(w http.ResponseWriter, r *http.Request) {
		staticPath := controllers.Page.Config.StaticDirectory + r.URL.Path[len(controllers.Page.Document.Static):]
		http.ServeFile(w, r, staticPath)
	})

	siteRootRightTrim := controllers.Page.Site.Root[:len(controllers.Page.Site.Root)-1]
	globalTpl := template.New("globalTpl").Funcs(controllers.Page.TemplateFunc)
	if t, _ := globalTpl.ParseGlob(controllers.Page.Config.TemplateDirectory + controllers.Page.Config.ThemeDirectory + controllers.Page.Config.TemplateGlobalDirectory + controllers.Page.Config.TemplateGlobalFile); t != nil {
		globalTpl = t
	}

	http.HandleFunc(controllers.Page.Site.Root, func(w http.ResponseWriter, r *http.Request) {
		controllers.Page.PageLock.Lock()
		yamlFi, _ := os.Stat(*configPath)
		if yamlFi.ModTime().Unix() > yamlLastModTime {
			yamlLastModTime = yamlFi.ModTime().Unix()
			yamlData, _ := ioutil.ReadFile(*configPath)
			controllers.Page.Config = controllers.Config
			goyaml.Unmarshal(yamlData, &controllers.Page.Config)

			controllers.Page.Site.Base.SupportSession = controllers.Page.Config.SupportSession
			controllers.Page.Config.UploadDirectory = controllers.Page.Config.StaticDirectory + controllers.Page.Config.UploadDirectory
			controllers.Page.Config.ThemeDirectory = controllers.Page.Config.ThemeDirectory + controllers.Page.Config.Theme + "/"
			controllers.Page.Config.StaticCssDirectory = controllers.Page.Config.StaticDirectory + controllers.Page.Config.ThemeDirectory + controllers.Page.Config.StaticCssDirectory
			controllers.Page.Config.StaticJsDirectory = controllers.Page.Config.StaticDirectory + controllers.Page.Config.ThemeDirectory + controllers.Page.Config.StaticJsDirectory
			controllers.Page.Config.StaticImgDirectory = controllers.Page.Config.StaticDirectory + controllers.Page.Config.ThemeDirectory + controllers.Page.Config.StaticImgDirectory

			if controllers.Page.Site.Root != controllers.Page.Config.SiteRoot {
				controllers.Page.UpdateController(controllers.Page.Site.Root, controllers.Page.Config.SiteRoot, controllers.Page.DefaultController)
				controllers.Page.Site.Root = controllers.Page.Config.SiteRoot
			}

			if t, _ := globalTpl.ParseGlob(controllers.Page.Config.TemplateDirectory + controllers.Page.Config.ThemeDirectory + controllers.Page.Config.TemplateGlobalDirectory + controllers.Page.Config.TemplateGlobalFile); t != nil {
				globalTpl = t
			}

			if globalCssFi, err := os.Stat(controllers.Page.Config.StaticCssDirectory + "/global/"); err == nil && globalCssFi.IsDir() {
				DcssPath := controllers.Page.Config.StaticCssDirectory + "global/"
				controllers.Page.Document.Css["global"] = siteRootRightTrim + DcssPath[1:]
				if _, err := os.Stat(DcssPath + "global.css"); err == nil {
					controllers.Page.Document.GlobalCssFile = controllers.Page.Document.Css["global"] + "global.css"
				}
			}

			if globalJsFi, err := os.Stat(controllers.Page.Config.StaticJsDirectory + "/global/"); err == nil && globalJsFi.IsDir() {
				DjsPath := controllers.Page.Config.StaticJsDirectory + "global/"
				controllers.Page.Document.Js["global"] = siteRootRightTrim + DjsPath[1:]
				if _, err := os.Stat(DjsPath + "global.js"); err == nil {
					controllers.Page.Document.GlobalJsFile = controllers.Page.Document.Js["global"] + "global.js"
				}
			}

			if globalImgFi, err := os.Stat(controllers.Page.Config.StaticImgDirectory + "/global/"); err == nil && globalImgFi.IsDir() {
				DimgPath := controllers.Page.Config.StaticImgDirectory + "global/"
				controllers.Page.Document.Img["global"] = siteRootRightTrim + DimgPath[1:]
			}
		}

		controllers.Page.Site.Base.Request = r
		controllers.Page.Site.Base.ResponseWriter = w
		controllers.Page.Site.Base.Cookie = r.Cookies()

		urlPath, fileName := filepath.Split(r.URL.Path)
		if urlPath == controllers.Page.Site.Root {
			urlPath = controllers.Page.Site.Root + controllers.Page.Config.IndexDirectory
		}

		if fileName == "" {
			fileName = controllers.Page.Config.IndexPage
		}

		controllers.Page.Template = urlPath[len(controllers.Page.Site.Root):] + fileName
		fileExt := filepath.Ext(fileName)
		fileNameNoExt := fileName[:len(fileName)-len(fileExt)]
		methodName := strings.Replace(strings.Title(strings.Replace(fileNameNoExt, "_", " ", -1)), " ", "", -1)

		if cssFi, err := os.Stat(controllers.Page.Config.StaticCssDirectory + urlPath); err == nil && cssFi.IsDir() {
			cssPath := strings.Trim(urlPath, "/")
			DcssPath := controllers.Page.Config.StaticCssDirectory + cssPath + "/"
			controllers.Page.Document.Css[cssPath] = siteRootRightTrim + DcssPath[1:]
			if _, err := os.Stat(DcssPath + "global.css"); err == nil {
				controllers.Page.Document.GlobalIndexCssFile = controllers.Page.Document.Css[cssPath] + "global.css"
			}

			if _, err := os.Stat(DcssPath + fileNameNoExt + ".css"); err == nil {
				controllers.Page.Document.IndexCssFile = controllers.Page.Document.Css[cssPath] + fileNameNoExt + ".css"
			}

		}

		if jsFi, err := os.Stat(controllers.Page.Config.StaticJsDirectory + urlPath); err == nil && jsFi.IsDir() {
			jsPath := strings.Trim(urlPath, "/")
			DjsPath := controllers.Page.Config.StaticJsDirectory + jsPath + "/"
			controllers.Page.Document.Js[jsPath] = siteRootRightTrim + DjsPath[1:]
			if _, err := os.Stat(DjsPath + "global.js"); err == nil {
				controllers.Page.Document.GlobalIndexJsFile = controllers.Page.Document.Js[jsPath] + "global.js"
			}

			if _, err := os.Stat(DjsPath + fileNameNoExt + ".js"); err == nil {
				controllers.Page.Document.IndexJsFile = controllers.Page.Document.Js[jsPath] + fileNameNoExt + ".js"
			}
		}

		if imgFi, err := os.Stat(controllers.Page.Config.StaticImgDirectory + urlPath); err == nil && imgFi.IsDir() {
			imgPath := strings.Trim(urlPath, "/")
			DimgPath := controllers.Page.Config.StaticImgDirectory + imgPath + "/"
			controllers.Page.Document.Img[imgPath] = siteRootRightTrim + DimgPath[1:]
		}

		if controllers.Page.Config.AutoGenerateHtml {
			controllers.Page.Document.GenerateHtml = true
		}

		controllers.Page.CurrentController = urlPath[len(controllers.Page.Site.Root):]
		controllers.Page.CurrentAction = methodName
		controllers.Page.PageLock.Unlock()

		controllers.Page.PageRLock.RLock()
		pageController := controllers.Page.GetController(controllers.Page.CurrentController)
		rv := reflect.ValueOf(pageController)
		rt := rv.Type()
		if _, ok := rt.MethodByName("Init"); ok {
			rv.MethodByName("Init").Call([]reflect.Value{})
		}

		if _, ok := rt.MethodByName(controllers.Page.CurrentAction); ok && controllers.Page.CurrentAction != "Init" && controllers.Page.Document.Close == false {
			rv.MethodByName(controllers.Page.CurrentAction).Call([]reflect.Value{})
		} else {
			if _, ok := pageController.(*controllers.Page404); !ok {
				controllers.Page.NotFoundtController.(*controllers.Page404).Init()
			}
		}

		if controllers.Page.Document.Close == false && controllers.Page.Document.Hide == false {
			if tplFi, err := os.Stat(controllers.Page.Config.TemplateDirectory + controllers.Page.Config.ThemeDirectory + controllers.Page.Template); err == nil {
				globalTemplate, _ := globalTpl.Clone()
				if pageTemplate, err := globalTemplate.New(filepath.Base(controllers.Page.Template)).ParseFiles(controllers.Page.Config.TemplateDirectory + controllers.Page.Config.ThemeDirectory + controllers.Page.Template); err == nil {
					templateVar := map[string]interface{}{
						"G":        controllers.Page.Base.GET,
						"P":        controllers.Page.Base.POST,
						"C":        controllers.Page.Base.COOKIE,
						"S":        controllers.Page.Base.SESSION,
						"Siteroot": controllers.Page.Site.Root,
						"Version":  controllers.Page.Site.Version,
						"Template": controllers.Page.Template,
						"D":        controllers.Page.Document,
						"Config":   controllers.Page.Config.M,
					}

					if controllers.Page.Document.GenerateHtml {
						htmlFile := controllers.Page.Config.StaticDirectory + controllers.Page.Config.HtmlDirectory + urlPath + fileName
						htmlDir := filepath.Dir(htmlFile)
						if htmlDirFi, err := os.Stat(htmlDir); err != nil || !htmlDirFi.IsDir() {
							os.MkdirAll(htmlDir, 0777)
						}

						var doWrite bool
						if controllers.Page.Config.AutoGenerateHtml {
							if controllers.Page.Config.AutoGenerateHtmlCycleTime <= 0 {
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
									case time.Now().Unix()-htmlFi.ModTime().Unix() >= controllers.Page.Config.AutoGenerateHtmlCycleTime:
										doWrite = true
									default:
										globalTplFi, err := os.Stat(controllers.Page.Config.TemplateDirectory + controllers.Page.Config.ThemeDirectory + controllers.Page.Config.TemplateGlobalDirectory)
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
								templateVar["Siteroot"] = controllers.Page.Config.SiteRoot + htmlDir + "/"
								controllers.Page.PageLock.Lock()
								pageTemplate.Execute(file, templateVar)
								controllers.Page.PageLock.Unlock()
							}
						}

						if controllers.Page.Config.AutoJumpToHtml {
							controllers.Page.PageLock.Lock()
							http.Redirect(w, r, controllers.Page.Site.Root+htmlFile[2:], http.StatusFound)
							controllers.Page.PageLock.Unlock()
						} else {
							controllers.Page.PageLock.Lock()
							err := pageTemplate.Execute(w, templateVar)
							controllers.Page.PageLock.Unlock()
							if err != nil {
								log.Println(err)
							}
						}
					} else {
						controllers.Page.PageLock.Lock()
						err := pageTemplate.Execute(w, templateVar)
						if err != nil {
							log.Println(err)
							w.Write([]byte(fmt.Sprint(err)))
						}

						controllers.Page.PageLock.Unlock()
					}
				} else {
					log.Println(err)
					controllers.Page.PageLock.Lock()
					w.Write([]byte(fmt.Sprint(err)))
					controllers.Page.PageLock.Unlock()
				}
			}
		}
		controllers.Page.PageRLock.RUnlock()

		controllers.Page.PageLock.Lock()
		controllers.Page.Reset()
		controllers.Page.PageLock.Unlock()
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
