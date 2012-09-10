package templateFunc

import (
	. "controllers"
	"golanger/utils"
)

func init() {
	App.AddTemplateFunc("GetTimeToStr", func(tm int64) string {
		return utils.NewTime().GetTimeToStr(tm)
	})
}
