package zillow

type MapBounds struct {
	West  float64 `json:"west"`
	East  float64 `json:"east"`
	South float64 `json:"south"`
	North float64 `json:"north"`
}

type SearchPageReq struct {
	MapBounds    *MapBounds `json:"mapBounds"`
	IsMapVisible bool       `json:"isMapVisible"`
	FilterState  struct {
		SortSelection struct {
			Value string `json:"value"`
		} `json:"sortSelection"`
		IsAllHomes struct {
			Value bool `json:"value"`
		} `json:"isAllHomes"`
	} `json:"filterState"`
	IsListVisible bool `json:"isListVisible"`
	MapZoom       int  `json:"mapZoom"`
	Pagination    struct {
		CurrentPage int `json:"currentPage"`
	} `json:"pagination"`
}
type SearchPageResResult struct {
	Zpid                  string      `json:"zpid"`
	ID                    string      `json:"id"`
	ProviderListingID     interface{} `json:"providerListingId"`
	StreetViewMetadataURL string      `json:"streetViewMetadataURL,omitempty"`
	StreetViewURL         string      `json:"streetViewURL,omitempty"`
	ImgSrc                string      `json:"imgSrc"`
	DetailURL             string      `json:"detailUrl"`
	StatusType            string      `json:"statusType"`
	StatusText            string      `json:"statusText"`
	CountryCurrency       string      `json:"countryCurrency"`
	Price                 string      `json:"price"`
	UnformattedPrice      float64     `json:"unformattedPrice"`
	Address               string      `json:"address"`
	AddressStreet         string      `json:"addressStreet"`
	AddressCity           string      `json:"addressCity"`
	AddressState          string      `json:"addressState"`
	AddressZipcode        string      `json:"addressZipcode"`
	IsUndisclosedAddress  bool        `json:"isUndisclosedAddress"`
	Beds                  float64     `json:"beds"`
	Baths                 float64     `json:"baths"`
	Area                  float64     `json:"area"`
	LatLong               struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"latLong"`
	IsZillowOwned bool `json:"isZillowOwned"`
	VariableData  struct {
		Type string `json:"type"`
		Text string `json:"text"`
		Data struct {
			IsFresh bool `json:"isFresh"`
		} `json:"data"`
	} `json:"variableData"`
	BadgeInfo interface{} `json:"badgeInfo"`
	HdpData   struct {
		HomeInfo struct {
			Zpid            int     `json:"zpid"`
			StreetAddress   string  `json:"streetAddress"`
			Zipcode         string  `json:"zipcode"`
			City            string  `json:"city"`
			State           string  `json:"state"`
			Latitude        float64 `json:"latitude"`
			Longitude       float64 `json:"longitude"`
			Price           float64 `json:"price"`
			Bathrooms       float64 `json:"bathrooms"`
			Bedrooms        float64 `json:"bedrooms"`
			LivingArea      float64 `json:"livingArea"`
			HomeType        string  `json:"homeType"`
			HomeStatus      string  `json:"homeStatus"`
			DaysOnZillow    float64 `json:"daysOnZillow"`
			IsFeatured      bool    `json:"isFeatured"`
			ShouldHighlight bool    `json:"shouldHighlight"`
			Zestimate       float64 `json:"zestimate"`
			RentZestimate   float64 `json:"rentZestimate"`
			ListingSubType  struct {
				IsFSBA bool `json:"is_FSBA"`
			} `json:"listing_sub_type"`
			IsUnmappable            bool    `json:"isUnmappable"`
			IsPreforeclosureAuction bool    `json:"isPreforeclosureAuction"`
			HomeStatusForHDP        string  `json:"homeStatusForHDP"`
			PriceForHDP             float64 `json:"priceForHDP"`
			IsNonOwnerOccupied      bool    `json:"isNonOwnerOccupied"`
			IsPremierBuilder        bool    `json:"isPremierBuilder"`
			IsZillowOwned           bool    `json:"isZillowOwned"`
			Currency                string  `json:"currency"`
			Country                 string  `json:"country"`
			TaxAssessedValue        float64 `json:"taxAssessedValue"`
			LotAreaValue            float64 `json:"lotAreaValue"`
			LotAreaUnit             string  `json:"lotAreaUnit"`
		} `json:"homeInfo"`
	} `json:"hdpData"`
	IsSaved                    bool        `json:"isSaved"`
	IsUserClaimingOwner        bool        `json:"isUserClaimingOwner"`
	IsUserConfirmedClaim       bool        `json:"isUserConfirmedClaim"`
	Pgapt                      string      `json:"pgapt"`
	Sgapt                      string      `json:"sgapt"`
	Zestimate                  float64     `json:"zestimate"`
	ShouldShowZestimateAsPrice bool        `json:"shouldShowZestimateAsPrice"`
	Has3DModel                 bool        `json:"has3DModel"`
	HasVideo                   bool        `json:"hasVideo"`
	IsHomeRec                  bool        `json:"isHomeRec"`
	Info2String                string      `json:"info2String,omitempty"`
	BrokerName                 string      `json:"brokerName,omitempty"`
	Info6String                string      `json:"info6String,omitempty"`
	HasAdditionalAttributions  bool        `json:"hasAdditionalAttributions"`
	IsFeaturedListing          bool        `json:"isFeaturedListing"`
	AvailabilityDate           interface{} `json:"availabilityDate"`
	List                       bool        `json:"list"`
	Relaxed                    bool        `json:"relaxed"`
	Info3String                string      `json:"info3String,omitempty"`
	HasImage                   bool        `json:"hasImage,omitempty"`
	BuilderName                string      `json:"builderName,omitempty"`
	Info1String                string      `json:"info1String,omitempty"`
	LotAreaString              string      `json:"lotAreaString,omitempty"`
}

type SearchPageResRelaxedResult struct {
	Zpid                 string      `json:"zpid"`
	ID                   string      `json:"id"`
	ProviderListingID    interface{} `json:"providerListingId"`
	ImgSrc               string      `json:"imgSrc"`
	HasImage             bool        `json:"hasImage"`
	DetailURL            string      `json:"detailUrl"`
	StatusType           string      `json:"statusType"`
	StatusText           string      `json:"statusText"`
	CountryCurrency      string      `json:"countryCurrency"`
	Price                string      `json:"price"`
	UnformattedPrice     float64     `json:"unformattedPrice"`
	Address              string      `json:"address"`
	AddressStreet        string      `json:"addressStreet"`
	AddressCity          string      `json:"addressCity"`
	AddressState         string      `json:"addressState"`
	AddressZipcode       string      `json:"addressZipcode"`
	IsUndisclosedAddress bool        `json:"isUndisclosedAddress"`
	Beds                 float64     `json:"beds"`
	Baths                float64     `json:"baths"`
	Area                 float64     `json:"area"`
	LatLong              struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"latLong"`
	IsZillowOwned bool `json:"isZillowOwned"`
	VariableData  struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"variableData"`
	BadgeInfo interface{} `json:"badgeInfo"`
	HdpData   struct {
		HomeInfo struct {
			Zpid            int     `json:"zpid"`
			StreetAddress   string  `json:"streetAddress"`
			Zipcode         string  `json:"zipcode"`
			City            string  `json:"city"`
			State           string  `json:"state"`
			Latitude        float64 `json:"latitude"`
			Longitude       float64 `json:"longitude"`
			Price           float64 `json:"price"`
			Bathrooms       float64 `json:"bathrooms"`
			Bedrooms        float64 `json:"bedrooms"`
			LivingArea      float64 `json:"livingArea"`
			HomeType        string  `json:"homeType"`
			HomeStatus      string  `json:"homeStatus"`
			DaysOnZillow    float64 `json:"daysOnZillow"`
			IsFeatured      bool    `json:"isFeatured"`
			ShouldHighlight bool    `json:"shouldHighlight"`
			Zestimate       float64 `json:"zestimate"`
			RentZestimate   float64 `json:"rentZestimate"`
			ListingSubType  struct {
				IsFSBA bool `json:"is_FSBA"`
			} `json:"listing_sub_type"`
			IsUnmappable            bool    `json:"isUnmappable"`
			IsPreforeclosureAuction bool    `json:"isPreforeclosureAuction"`
			HomeStatusForHDP        string  `json:"homeStatusForHDP"`
			PriceForHDP             float64 `json:"priceForHDP"`
			IsNonOwnerOccupied      bool    `json:"isNonOwnerOccupied"`
			IsPremierBuilder        bool    `json:"isPremierBuilder"`
			IsZillowOwned           bool    `json:"isZillowOwned"`
			Currency                string  `json:"currency"`
			Country                 string  `json:"country"`
			TaxAssessedValue        float64 `json:"taxAssessedValue"`
			LotAreaValue            float64 `json:"lotAreaValue"`
			LotAreaUnit             string  `json:"lotAreaUnit"`
		} `json:"homeInfo"`
	} `json:"hdpData"`
	IsSaved                    bool        `json:"isSaved"`
	IsUserClaimingOwner        bool        `json:"isUserClaimingOwner"`
	IsUserConfirmedClaim       bool        `json:"isUserConfirmedClaim"`
	Pgapt                      string      `json:"pgapt"`
	Sgapt                      string      `json:"sgapt"`
	Zestimate                  float64     `json:"zestimate"`
	ShouldShowZestimateAsPrice bool        `json:"shouldShowZestimateAsPrice"`
	Has3DModel                 bool        `json:"has3DModel"`
	HasVideo                   bool        `json:"hasVideo"`
	IsHomeRec                  bool        `json:"isHomeRec"`
	HasAdditionalAttributions  bool        `json:"hasAdditionalAttributions"`
	IsFeaturedListing          bool        `json:"isFeaturedListing"`
	AvailabilityDate           interface{} `json:"availabilityDate"`
	Relaxed                    bool        `json:"relaxed"`
	Ma                         bool        `json:"ma"`
	HasOpenHouse               bool        `json:"hasOpenHouse,omitempty"`
	OpenHouseStartDate         string      `json:"openHouseStartDate,omitempty"`
	OpenHouseEndDate           string      `json:"openHouseEndDate,omitempty"`
	OpenHouseDescription       string      `json:"openHouseDescription,omitempty"`
	BrokerName                 string      `json:"brokerName,omitempty"`
}

type SearchPageRes struct {
	User struct {
		IsLoggedIn                    bool        `json:"isLoggedIn"`
		HasHousingConnectorPermission bool        `json:"hasHousingConnectorPermission"`
		SavedSearchCount              int         `json:"savedSearchCount"`
		SavedHomesCount               int         `json:"savedHomesCount"`
		PersonalizedSearchGaDataTag   interface{} `json:"personalizedSearchGaDataTag"`
		PersonalizedSearchTraceID     string      `json:"personalizedSearchTraceID"`
		SearchPageRenderedCount       int         `json:"searchPageRenderedCount"`
		GUID                          string      `json:"guid"`
		Zuid                          string      `json:"zuid"`
		IsBot                         bool        `json:"isBot"`
		UserSpecializedSEORegion      bool        `json:"userSpecializedSEORegion"`
	} `json:"user"`
	MapState struct {
		CustomRegionPolygonWkt  interface{} `json:"customRegionPolygonWkt"`
		SchoolPolygonWkt        interface{} `json:"schoolPolygonWkt"`
		IsCurrentLocationSearch bool        `json:"isCurrentLocationSearch"`
		UserPosition            struct {
			Lat interface{} `json:"lat"`
			Lon interface{} `json:"lon"`
		} `json:"userPosition"`
		RegionInfo []interface{} `json:"regionInfo"`
	} `json:"mapState"`
	SearchPageSeoObject struct {
		BaseURL         string `json:"baseUrl"`
		WindowTitle     string `json:"windowTitle"`
		MetaDescription string `json:"metaDescription"`
	} `json:"searchPageSeoObject"`
	RequestID int `json:"requestId"`
	Cat1      struct {
		SearchResults struct {
			ListResults        []SearchPageResResult        `json:"listResults"`
			ResultsHash        string                       `json:"resultsHash"`
			HomeRecCount       int                          `json:"homeRecCount"`
			ShowForYouCount    int                          `json:"showForYouCount"`
			MapResults         []interface{}                `json:"mapResults"`
			RelaxedResults     []SearchPageResRelaxedResult `json:"relaxedResults"`
			RelaxedResultsHash string                       `json:"relaxedResultsHash"`
		} `json:"searchResults"`
		SearchList struct {
			ExpansionDistance  int         `json:"expansionDistance"`
			ZeroResultsFilters interface{} `json:"zeroResultsFilters"`
			Pagination         struct {
				NextURL string `json:"nextUrl"`
			} `json:"pagination"`
			Message   interface{} `json:"message"`
			AdsConfig struct {
				NavAdSlot     string `json:"navAdSlot"`
				DisplayAdSlot string `json:"displayAdSlot"`
				Targets       struct {
					Mlat         string `json:"mlat"`
					Zusr         string `json:"zusr"`
					Listtp       string `json:"listtp"`
					Searchtp     string `json:"searchtp"`
					Filtered     string `json:"filtered"`
					Premieragent string `json:"premieragent"`
					GUID         string `json:"guid"`
					Mlong        string `json:"mlong"`
				} `json:"targets"`
				NeedsUpdate bool `json:"needsUpdate"`
			} `json:"adsConfig"`
			TotalResultCount        int    `json:"totalResultCount"`
			ResultsPerPage          int    `json:"resultsPerPage"`
			TotalPages              int    `json:"totalPages"`
			LimitSearchResultsCount int    `json:"limitSearchResultsCount"`
			ListResultsTitle        string `json:"listResultsTitle"`
			ResultContexts          []struct {
				Ssid         int    `json:"ssid"`
				Context      string `json:"context"`
				ContextImage string `json:"contextImage"`
			} `json:"resultContexts"`
			PageRules string `json:"pageRules"`
		} `json:"searchList"`
	} `json:"cat1"`
	CategoryTotals struct {
		Cat1 struct {
			TotalResultCount int `json:"totalResultCount"`
		} `json:"cat1"`
		Cat2 struct {
			TotalResultCount int `json:"totalResultCount"`
		} `json:"cat2"`
	} `json:"categoryTotals"`
}

type ZillowData struct {
	URL                        string
	Address                    string
	PropertyStatus             bool
	Bed                        int
	Bath                       float64
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
	LotSizeSF                  string
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
