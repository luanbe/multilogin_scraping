package schemas

import (
	"errors"
	"multilogin_scraping/app/models/entity"
	"net/http"
)

type CrawlerRequest struct {
	Address string `json:"address"`
}

type CrawlerResponse struct {
	TaskID  string `json:"task_id"`
	Address string `json:"address"`
}

type CrawlerTask struct {
	Status       string               `json:"status"`
	TaskID       string               `json:"task_id"`
	Address      string               `json:"address"`
	Error        string               `json:"error"`
	ZillowDetail *entity.ZillowDetail `json:"zillow_detail"`
}

func (c *CrawlerRequest) Bind(r *http.Request) error {
	if c.Address == "" {
		return errors.New("missing required Address fields.")
	}
	return nil
}
func (cr *CrawlerResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}
func (cr *CrawlerTask) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func NewCrawlerResponse(taskID, address string) *CrawlerResponse {
	resp := &CrawlerResponse{taskID, address}
	return resp
}
