package controllers

type PageIndex struct {
	*Application
}

func init() {
	App.RegisterController("index/", &PageIndex{App})
}

func (p *PageIndex) Index() {
}
