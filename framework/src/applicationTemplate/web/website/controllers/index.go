package controllers

type PageIndex struct {
	*App
}

func init() {
	Page.RegisterController("index/", &PageIndex{Page})
}

func (p *PageIndex) Index() {
}
