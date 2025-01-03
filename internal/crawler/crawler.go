package crawler

import (
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/sourcegraph/conc/pool"

	. "adc/internal/config"
	. "adc/internal/javInfo"
	. "adc/internal/logger"
	. "adc/internal/videos"
)

type AvInfoCrawler interface {
	Init()
	CrawlAdultVideo()
}

func StartTasks(mode string) {
	if mode != "one-time" && mode != "watchdog" {
		Logger.Errorf("mode %s is not valid!", mode)
		return
	}

	if err := JavWalk(AdcConfig.SourceDirectory); err != nil {
		Logger.Errorln(err)
		return
	}

	if mode == "one-time" {
		Logger.Infof("Found %d videos under %s", len(FilesList), AdcConfig.SourceDirectory)
	}

	startCrawlers(FilesList, true)
}

func startCrawlers(fileLists []*JavInfo, enableOrganize bool) {
	p := pool.New().WithMaxGoroutines(5)

	var opts []JavbusOptions
	if enableOrganize {
		opts = append(opts, WithOrganize(true))
	}

	for _, info := range fileLists {
		var crawler AvInfoCrawler
		switch info.Type {
		case "jav":
			crawler = NewJavbusCrawler(info, opts...)
		default:
			Logger.WithField("number", info.Number).Errorf("%s is not supported yet!", info.Type)
			continue
		}
		crawler.Init()
		p.Go(crawler.CrawlAdultVideo)
	}
	p.Wait()
}

func newGeneralInfoCollector() *colly.Collector {
	collector := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(),
	)

	// dev mode
	logLevel := AdcConfig.LoggerOptions.Level
	if logLevel == "dev" {
		collector.SetDebugger(&CrawlLogger{})
		collector.Async = false
	}

	// rate limit
	err := collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})
	if err != nil {
		Logger.Errorln(err)
	}

	timeout := time.Duration(AdcConfig.Proxy.Timeout) * time.Second
	collector.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: timeout,
		}).DialContext,
	})

	// anti-anti-crawler options
	extensions.RandomUserAgent(collector)
	extensions.Referer(collector)

	// proxy balancer
	if rp, err := proxy.RoundRobinProxySwitcher(AdcConfig.Proxy.URL); err != nil {
		Logger.Errorln(err)
	} else {
		collector.SetProxyFunc(rp)
	}

	return collector
}
