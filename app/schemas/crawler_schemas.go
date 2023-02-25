package schemas

import (
	"errors"
	"multilogin_scraping/app/models/entity"
	"net/http"
)

type CrawlerRequest struct {
	Search         CrawlerSearch  `json:"search"`
	CrawlersActive CrawlersStatus `json:"crawlers_status"`
}

type CrawlerSearch struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zipcode string `json:"zipcode"`
}
type CrawlersStatus struct {
	Zillow  bool `json:"zillow"`
	Realtor bool `json:"realtor"`
	Movoto  bool `json:"movoto"`
}

func (c *CrawlerRequest) Bind(r *http.Request) error {
	if c.Search.Address == "" || c.Search.City == "" || c.Search.State == "" || c.Search.Zipcode == "" {
		return errors.New("missing required address, city, state and zipcode field.")
	}
	return nil
}

type CrawlerSearchRes struct {
	TaskID         string          `json:"task_id"`
	CrawlerRequest *CrawlerRequest `json:"crawler_request"`
	Zillow         *ZillowDetail   `json:"zillow"`
	Realtor        *RealtorDetail  `json:"realtor"`
	Movoto         *MovotoDetail   `json:"movoto"`
}

type ZillowDetail struct {
	Status string         `json:"status"`
	Error  string         `json:"error"`
	Data   *entity.Zillow `json:"data"`
}

type RealtorDetail struct {
	Status string          `json:"status"`
	Error  string          `json:"error"`
	Data   *entity.Realtor `json:"data"`
}

type MovotoDetail struct {
	Status string         `json:"status"`
	Error  string         `json:"error"`
	Data   *entity.Movoto `json:"data"`
}

func (ct *CrawlerSearchRes) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func NewCrawlerResponse(taskID string, crawlerSearchRes CrawlerSearchRes) *CrawlerSearchRes {
	return &crawlerSearchRes
}
