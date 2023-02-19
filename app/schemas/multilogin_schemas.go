package schemas

type MultiLoginProfile []struct {
	UUID               string      `json:"uuid"`
	Group              string      `json:"group"`
	Name               string      `json:"name"`
	Notes              string      `json:"notes"`
	Browser            string      `json:"browser"`
	BrowserNeedsUpdate interface{} `json:"browserNeedsUpdate"`
	AppNeedsUpdate     interface{} `json:"appNeedsUpdate"`
	Updated            float64     `json:"updated"`
}
