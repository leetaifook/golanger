package controllers

type PagePlay struct {
	*App
}

func init() {
	Page.RegisterController("play/", &PagePlay{Page})
}

func (p *PagePlay) Index() {
}
