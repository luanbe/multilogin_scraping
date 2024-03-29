package schemas

import (
	"net/http"
)

type RealtorSearchPageReq struct {
	Input     string `json:"input" url:"input"`
	ClientID  string `json:"client_id" url:"client_id"`
	Limit     int    `json:"limit" url:"limit"`
	AreaTypes string `json:"area_types" url:"area_types"`
}
type RealtorSearchPageRes struct {
	Meta struct {
		Version string `json:"version"`
		EsTook  int    `json:"es_took"`
	} `json:"meta"`
	Autocomplete []struct {
		AreaType    string   `json:"area_type"`
		ID          string   `json:"_id"`
		Score       float64  `json:"_score"`
		MprID       string   `json:"mpr_id"`
		FullAddress []string `json:"full_address"`
		Line        string   `json:"line"`
		City        string   `json:"city"`
		PostalCode  string   `json:"postal_code"`
		StateCode   string   `json:"state_code"`
		Country     string   `json:"country"`
		Centroid    struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		} `json:"centroid"`
		PropStatus     []string `json:"prop_status"`
		ValidationCode []string `json:"validation_code"`
	} `json:"autocomplete"`
}

type RealtorData struct {
	URL                        string
	Address                    string
	PropertyStatus             bool
	Bed                        int
	Bath                       int
	FullBathrooms              float64
	HalfBathrooms              float64
	SF                         float64
	SalesPrice                 float64
	EstPayment                 string
	PrincipalInterest          string
	MortgageInsurance          string
	PropertyTaxes              string
	HomeInsurance              string
	HOAFee                     string
	Utilities                  string
	RentZestimate              float64
	Zestimate                  float64
	EstimatedSalesRangeMinimum string
	EstimatedSalesRangeMax     string
	Pictures                   []string
	TimeOnZillow               string
	Views                      int
	Saves                      int
	Overview                   string
	ZillowCheckedDate          string
	DataUploadedDate           string
	ListedBy                   []string
	Source                     string
	MLS                        string
	PropertyType               string
	YearBuilt                  string
	NaturalGas                 bool
	CentralAir                 bool
	OfGarageSpaces             string
	HOAAmount                  string
	LotSizeSF                  float64
	LotSizeAcres               string
	BuyerAgentFee              string
	Appliances                 string
	LivingRoomLevel            string
	LivingRoomDimensions       string
	InteriorFeatures           string
	PrimaryBedroomLevel        string
	PrimaryBedroomDimensions   string
	Basement                   string
	TotalInteriorLivableAreaSF string
	OfFireplaces               string
	FireplaceFeatures          string
	FlooringType               string
	HeatingType                string
	TotalParkingSpaces         string
	ParkingFeatures            string
	LotFeatures                string
	CoveredSpaces              string
	ParcelNumber               string
	LevelsStoriesFloors        string
	PatioAndPorchDetails       string
	HomeType                   string
	ProperySubType             string
	ConstructionMaterials      string
	Foundation                 string
	Roof                       string
	NewConstruction            string
	SewerInformation           string
	WaterInformation           string
	RegionLocation             string
	Subdivision                string
	HasHOA                     string
	HOAFeeDetail               string
	ServicesIncluded           string
	AssociationName            string
	AssociationPhone           string
	AnnualTaxAmount            string
	ElementarySchool           string
	MiddleSchool               string
	HighSchool                 string
	District                   string
	DataSource                 string
	CountyTaxAssessorURL       string
	TimestampForDataExtraction string
}

type RealtorCrawlerTask struct {
	Status        string       `json:"status"`
	TaskID        string       `json:"task_id"`
	Address       string       `json:"address"`
	Error         string       `json:"error"`
	RealtorDetail *RealtorData `json:"realtor_detail"`
}

func (rc *RealtorCrawlerTask) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}
