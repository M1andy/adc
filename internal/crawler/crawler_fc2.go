package crawler

import (
	"github.com/gocolly/colly/v2"

	. "adc/internal/javInfo"
)

var (
	fc2HomePage = "https://adult.contents.fc2.com/"
	moviePage   = "https://adult.contents.fc2.com/article/%s"
)

type FC2Crawler struct {
	info      *JavInfo
	collector *colly.Collector
}

func NewFC2Crawler(info *JavInfo) *FC2Crawler {
	// init default collector
	collector := newGeneralInfoCollector()

	// setup domains
	collector.AllowedDomains = javbusDomains

	crawler := &FC2Crawler{
		info:      info,
		collector: collector,
	}

	return crawler
}

func (c FC2Crawler) Init() {
	//TODO implement me
	panic("implement me")
}

func (c FC2Crawler) CrawlAdultVideo() {
	//TODO implement me
	panic("implement me")
}
