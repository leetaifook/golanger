package helper

import (
	. "../controllers"
	"golanger/utils"
)

func init() {
	Page.AddTemplateFunc("GetTimeToStr", func(tm int64) string {
		return utils.NewTime().GetTimeToStr(tm)
	})
}
