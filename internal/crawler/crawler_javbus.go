package crawler

import (
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	str "strings"

	"github.com/gocolly/colly/v2"

	"adc/internal/videos"

	. "adc/internal/config"
	. "adc/internal/javInfo"
	. "adc/internal/logger"
)

var (
	javbusHomePage = "www.javbus.com"
	javbusDomains  = []string{
		// javbus websites
		"www.javbus.com",
		"www.javsee.help",
		"www.fanbus.help",

		// preview pics website
		"pics.dmm.co.jp",
		"pics-cache-digcdp.dmm.com",
	}
)

var (
	removeEmptyCharsReg = regexp.MustCompile(`[\n\t ]+`)
)

type JavbusCrawler struct {
	info       *JavInfo
	collector  *colly.Collector
	isOrganize bool
}

type JavbusOptions func(*JavbusCrawler)

func WithOrganize(isOrganize bool) JavbusOptions {
	return func(j *JavbusCrawler) {
		j.isOrganize = isOrganize
	}
}

func NewJavbusCrawler(info *JavInfo, opts ...JavbusOptions) *JavbusCrawler {
	// init default collector
	collector := newGeneralInfoCollector()

	// setup domains
	collector.AllowedDomains = javbusDomains

	crawler := &JavbusCrawler{
		info:      info,
		collector: collector,
	}

	for _, opt := range opts {
		opt(crawler)
	}
	return crawler
}

func (c *JavbusCrawler) Init() {
	// setup header
	c.collector.OnRequest(func(r *colly.Request) {
		// other anti-anti-crawler headers
		r.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		r.Headers.Set("accept-language", "zh-CN,zh;q=0.9")
		r.Headers.Set("cache-control", "max-age=0")
		r.Headers.Set("priority", "u=0, i")
		r.Headers.Set("sec-fetch-dest", "document")
		r.Headers.Set("sec-fetch-mode", "navigate")
		r.Headers.Set("sec-fetch-site", "none")
		r.Headers.Set("sec-fetch-user", "?1")
		r.Headers.Set("upgrade-insecure-requests", "1")
	})

	// title
	titleSrcPath := "//div[@class='container']"
	c.collector.OnXML(titleSrcPath, func(e *colly.XMLElement) {
		title := e.ChildText("h3/text()")
		titleWithoutNumber := str.TrimLeft(str.ToUpper(title), c.info.Number)
		c.info.Title = str.TrimSpace(titleWithoutNumber)
	})

	// general info
	generalInfoSrcPath := "//div[@class='col-md-3 info']"
	c.collector.OnXML(generalInfoSrcPath, func(e *colly.XMLElement) {
		number := c.info.Number
		infos := e.ChildTexts("p")

		javLogger := Logger.WithField("number", number)

		for i := range infos {
			info := infos[i]
			kv := str.Split(info, ":")

			if len(kv) == 1 {
				continue
			}

			key := kv[0]
			value := str.TrimSpace(kv[1])

			switch key {
			case "識別碼":
				if number != value {
					javLogger.Warnf("CrawlNumber Wrong! SrcNumber: %s | CrawlNumber: %s", number, value)
					break
				}
			case "發行日期":
				c.info.ReleaseDate = value
			case "長度":
				c.info.VideoLength = str.ReplaceAll(value, "分鐘", "分钟")
			case "製作商":
				c.info.Manufacturer = value
			case "發行商":
				c.info.Studio = value
			case "系列":
				c.info.Series = value
			case "類別":
				c.info.Genre = splitEmptySepStrings(infos[i+1])
			case "演員":
				c.info.Actresses = splitEmptySepStrings(infos[i+1])
			default:
				continue
			}
		}

		// update outPath
		outPath := AdcConfig.SuccessOutputDirectory
		actress, err := parseActressesList(c.info.Actresses)
		if err != nil {
			javLogger.Warnln(err)
			actress = "未知演员"
		}
		infoOutPath := path.Join(outPath, actress, number)
		c.info.OutDir = infoOutPath

		//debug info
		j, _ := json.MarshalIndent(c.info, "", "\t")
		javLogger.Debugln(string(j))
	})
}

func (c *JavbusCrawler) CrawlAdultVideo() {
	infoUrl := fmt.Sprintf("https://%s/%s", javbusHomePage, c.info.Number)
	err := c.collector.Visit(infoUrl)
	c.collector.Wait()

	number := c.info.Number
	javLogger := Logger.WithField("number", number)

	if err != nil {
		javLogger.Errorf("Crawl info failed! %s", err)
		return
	}

	if c.isOrganize {
		organizeJav(c.info)
	}
}

func organizeJav(info *JavInfo) {
	number := info.Number
	javLogger := Logger.WithField("number", number)

	srcFilePath := info.SrcFilePath
	outDir := info.OutDir
	err := videos.MoveJav(srcFilePath, outDir)
	if err != nil {
		javLogger.Error(err)
	}

	javLogger.Infof("JavFile move to %s", outDir)
}

func splitEmptySepStrings(s string) []string {
	actorClean := removeEmptyCharsReg.ReplaceAllString(s, ":")
	actors := str.Split(actorClean, ":")
	return actors
}

func parseActressesList(actresses []string) (string, error) {
	if len(actresses) == 0 {
		return "", fmt.Errorf("actress not found")
	}

	if len(actresses) == 1 {
		return actresses[0], nil
	}

	maxLen := min(3, len(actresses))

	return str.Join(actresses[:maxLen], ","), nil
}
