package main

import (
	"multilogin_scraping/crawlers"
)

func main() {
	zillow := crawlers.NewZillowCrawler()
	zillow.RunZillowCrawler()
}
