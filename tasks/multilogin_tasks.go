package tasks

import (
	"go.uber.org/zap"
	"multilogin_scraping/crawlers"
)

type MultiLoginProcessor struct {
	Logger *zap.Logger
}

func (m MultiLoginProcessor) DeleteProfiles(profiles []string) {
	baseSel := crawlers.NewBaseSelenium(m.Logger)
	err := baseSel.Profile.DeleteProfiles(profiles, m.Logger)
	if err != nil {
		m.Logger.Error(err.Error())
	}
}
