package web

type Document struct {
	GenerateHtml       bool
	Static             string
	Theme              string
	Css                map[string]string
	Js                 map[string]string
	Img                map[string]string
	GlobalCssFile      string
	GlobalJsFile       string
	GlobalIndexCssFile string
	GlobalIndexJsFile  string
	IndexCssFile       string
	IndexJsFile        string
	Hide               bool
	Title              string
	Subtitle           string
	Header             string
	Body               interface{}
	Footer             string
}
