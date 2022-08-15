package main

import (
	"github.com/gocolly/colly"
	"multilogin_scraping/crawlers/zillow"
)

func main() {
	c := colly.NewCollector()
	cZillow := c.Clone()
	zillow := zillow.NewZillowCrawler(cZillow)
	zillow.RunZillowCrawler()
}
