package schemas

import (
	"net/http"
)

type MovotoSearchPageReq struct {
	Input     string `json:"input" url:"input"`
	ClientID  string `json:"client_id" url:"client_id"`
	Limit     int    `json:"limit" url:"limit"`
	AreaTypes string `json:"area_types" url:"area_types"`
}

type MovotoSearchPageRes struct {
	Data struct {
		SearchCondition struct {
			Input                 string        `json:"input"`
			PropertyTypes         []interface{} `json:"propertyTypes"`
			PageIndex             int           `json:"pageIndex"`
			MaxCountPerPage       int           `json:"maxCountPerPage"`
			IsReducedPrice        int           `json:"isReducedPrice"`
			IsOpenHousesOnly      int           `json:"isOpenHousesOnly"`
			IsVirtualTourLinkOnly int           `json:"isVirtualTourLinkOnly"`
			IsNewConstruction     int           `json:"isNewConstruction"`
			IsNewListingsOnly     int           `json:"isNewListingsOnly"`
			IsRental              int           `json:"isRental"`
			AttributesTags        []interface{} `json:"attributesTags"`
			IsDistressed          int           `json:"isDistressed"`
			IsFixerUpper          int           `json:"isFixerUpper"`
			HasPhoto              int           `json:"hasPhoto"`
			HasPool               int           `json:"hasPool"`
			MinLat                int           `json:"minLat"`
			MaxLat                int           `json:"maxLat"`
			MinLng                int           `json:"minLng"`
			MaxLng                int           `json:"maxLng"`
			SearchPropertyStatus  string        `json:"searchPropertyStatus"`
			SchoolDistricts       interface{}   `json:"schoolDistricts"`
			SearchType            string        `json:"searchType"`
			MapView               bool          `json:"mapView"`
			LuxuryHomes           bool          `json:"luxuryHomes"`
			ExpandListingAmount   int           `json:"expandListingAmount"`
			DefaultSearch         bool          `json:"defaultSearch"`
			MovotoListing         int           `json:"movotoListing"`
			SelfTourable          bool          `json:"selfTourable"`
			Log                   string        `json:"log"`
			InputDisplay          string        `json:"inputDisplay"`
			SortOrder             string        `json:"sortOrder"`
			SortColumn            string        `json:"sortColumn"`
			IsHomeRoam            bool          `json:"isHomeRoam"`
			NonCompliance         bool          `json:"nonCompliance"`
		} `json:"searchCondition"`
		ViewURL    string   `json:"viewUrl"`
		Attributes []string `json:"attributes"`
		Listings   []struct {
			ClosePrice   interface{} `json:"closePrice"`
			DaysOnMovoto int         `json:"daysOnMovoto"`
			ID           string      `json:"id"`
			TnImgPath    string      `json:"tnImgPath"`
			ListDate     interface{} `json:"listDate"`
			ListingAgent string      `json:"listingAgent"`
			ListPrice    int         `json:"listPrice"`
			LotSize      int         `json:"lotSize"`
			SqftTotal    int         `json:"sqftTotal"`
			MlsDbNumber  int         `json:"mlsDbNumber"`
			Mls          struct {
				DateHidden interface{} `json:"dateHidden"`
				ID         int         `json:"id"`
			} `json:"mls"`
			MlsNumber               string        `json:"mlsNumber"`
			Bath                    int           `json:"bath"`
			Bed                     int           `json:"bed"`
			OpenHouses              []interface{} `json:"openHouses"`
			OfficeListName          string        `json:"officeListName"`
			PhotoCount              int           `json:"photoCount"`
			PropertyType            string        `json:"propertyType"`
			PropertyTypeDisplayName string        `json:"propertyTypeDisplayName"`
			YearBuilt               int           `json:"yearBuilt"`
			ZipCode                 string        `json:"zipCode"`
			Path                    string        `json:"path"`
			Status                  string        `json:"status"`
			HouseRealStatus         string        `json:"houseRealStatus"`
			Hoafee                  int           `json:"hoafee"`
			PriceChangedDate        string        `json:"priceChangedDate"`
			PriceChange             int           `json:"priceChange"`
			PropertyID              string        `json:"propertyId"`
			Visibility              string        `json:"visibility"`
			SoldDate                interface{}   `json:"soldDate"`
			CreatedAt               string        `json:"createdAt"`
			ImageDownloaderStatus   int           `json:"imageDownloaderStatus"`
			OnMarketDateTime        string        `json:"onMarketDateTime"`
			VirtualTourLink         string        `json:"virtualTourLink"`
			NhsDetails              interface{}   `json:"nhsDetails"`
			RentalDetails           interface{}   `json:"rentalDetails"`
			BuildingName            interface{}   `json:"buildingName"`
			PropertySizeSort        int           `json:"propertySizeSort"`
			BrokerageDetails        interface{}   `json:"brokerageDetails"`
			Geo                     struct {
				State            string      `json:"state"`
				City             string      `json:"city"`
				Lat              float64     `json:"lat"`
				Lng              float64     `json:"lng"`
				Zipcode          string      `json:"zipcode"`
				SubPremise       string      `json:"subPremise"`
				Address          string      `json:"address"`
				NeighborhoodName interface{} `json:"neighborhoodName"`
			} `json:"geo"`
			IsNHS           bool   `json:"isNHS"`
			IsRentals       bool   `json:"isRentals"`
			IsSold          bool   `json:"isSold"`
			ListingByMovoto bool   `json:"listingByMovoto"`
			PriceRaw        int    `json:"priceRaw"`
			IsVideoTour     bool   `json:"isVideoTour"`
			Is3DTour        bool   `json:"is3DTour"`
			VideoTourLink   string `json:"videoTourLink"`
			VirtualLink     string `json:"virtualLink"`
		} `json:"listings"`
		MlsIds          []int       `json:"mlsIds"`
		BoundaryIndexID interface{} `json:"boundaryIndexId"`
		TotalCount      int         `json:"totalCount"`
		SearchType      string      `json:"searchType"`
		ListingCountObj struct {
		} `json:"listingCountObj"`
		TopZipCodeList      interface{} `json:"topZipCodeList"`
		TopNearbyCityList   interface{} `json:"topNearbyCityList"`
		TopNeighborhoodList interface{} `json:"topNeighborhoodList"`
		TopNearbyCountyList interface{} `json:"topNearbyCountyList"`
		GeoPhone            string      `json:"geoPhone"`
		SearchText          string      `json:"searchText"`
		SchoolData          interface{} `json:"schoolData"`
		PropertyTypeEnum    struct {
			SINGLEFAMILYHOUSE string `json:"SINGLE_FAMILY_HOUSE"`
			CONDO             string `json:"CONDO"`
			MULTIFAMILY       string `json:"MULTI_FAMILY"`
			LAND              string `json:"LAND"`
			COMMERCIAL        string `json:"COMMERCIAL"`
			MOBILE            string `json:"MOBILE"`
			OTHER             string `json:"OTHER"`
		} `json:"propertyTypeEnum"`
		Seotitle struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			H1          string `json:"h1"`
		} `json:"seotitle"`
		FilterType struct {
		} `json:"filterType"`
		CurrentPageIndex int `json:"currentPageIndex"`
		PageSize         int `json:"pageSize"`
		RelatedLinks     struct {
			Show bool `json:"show"`
		} `json:"relatedLinks"`
		HasFilter       bool `json:"hasFilter"`
		HasAttr         bool `json:"hasAttr"`
		ShowFilterLinks bool `json:"showFilterLinks"`
		ShowAttrLinks   bool `json:"showAttrLinks"`
		AttrSSRLinks    struct {
			Pool           string `json:"pool"`
			Backyard       string `json:"backyard"`
			OpenFloorPlan  string `json:"open_floor_plan"`
			WalkInCloset   string `json:"walk_in_closet"`
			Fireplace      string `json:"fireplace"`
			Clubhouse      string `json:"clubhouse"`
			Patio          string `json:"patio"`
			Deck           string `json:"deck"`
			SingleLevel    string `json:"single_level"`
			HardwoodFloors string `json:"hardwood_floors"`
			View           string `json:"view"`
			LargeKitchen   string `json:"large_kitchen"`
			KitchenIsland  string `json:"kitchen_island"`
			Porch          string `json:"porch"`
			Garage         string `json:"garage"`
			Waterfront     string `json:"waterfront"`
		} `json:"attrSSRLinks"`
	} `json:"data"`
	Status struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"status"`
	Version string `json:"version"`
}

type MovotoData struct {
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

type MovotoCrawlerTask struct {
	Status       string      `json:"status"`
	TaskID       string      `json:"task_id"`
	Address      string      `json:"address"`
	Error        string      `json:"error"`
	MovotoDetail *MovotoData `json:"realtor_detail"`
}

func (rc *MovotoCrawlerTask) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}
