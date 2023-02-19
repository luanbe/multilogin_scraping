package schemas

import (
	"errors"
	"multilogin_scraping/app/models/entity"
	"net/http"
)

type CrawlerRequest struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zipcode string `json:"zipcode"`
}

func (c *CrawlerRequest) Bind(r *http.Request) error {
	if c.Address == "" {
		return errors.New("missing required Address fields.")
	}
	return nil
}

type CrawlerSearchRes struct {
	TaskID  string          `json:"task_id"`
	Search  *CrawlerRequest `json:"search"`
	Zillow  *ZillowDetail   `json:"zillow"`
	Realtor *RealtorDetail  `json:"realtor"`
	Movoto  *MovotoDetail   `json:"movoto"`
}

type ZillowDetail struct {
	Status string               `json:"status"`
	Error  string               `json:"error"`
	Data   *entity.ZillowDetail `json:"data"`
}

type RealtorDetail struct {
	Status string       `json:"status"`
	Error  string       `json:"error"`
	Data   *RealtorData `json:"data"`
}

type MovotoDetail struct {
	Status string      `json:"status"`
	Error  string      `json:"error"`
	Data   *MovotoData `json:"data"`
}

func (ct *CrawlerSearchRes) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func NewCrawlerResponse(taskID string, crawlerSearchRes CrawlerSearchRes) *CrawlerSearchRes {
	return &crawlerSearchRes
}
