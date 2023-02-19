package schemas

import (
	"net/http"
	"time"
)

type MovotoSearchData struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zipcode string `json:"zipcode"`
}
type MovotoSearchPageReq struct {
	Path              string `json:"path" url:"path"`
	Trigger           string `json:"trigger" url:"trigger"`
	IncludeAllAddress bool   `json:"includeAllAddress" url:"includeAllAddress"`
	NewGeoSearch      bool   `json:"newGeoSearch" url:"newGeoSearch"`
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
		ViewURL         string                `json:"viewUrl"`
		Attributes      []string              `json:"attributes"`
		Listings        []MovotoSearchDataRes `json:"listings"`
		MlsIds          []int                 `json:"mlsIds"`
		BoundaryIndexID interface{}           `json:"boundaryIndexId"`
		TotalCount      int                   `json:"totalCount"`
		SearchType      string                `json:"searchType"`
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

type MovotoSearchDataRes struct {
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
	IsNHS     bool `json:"isNHS"`
	IsRentals bool `json:"isRentals"`
	IsSold    bool `json:"isSold"`
	//ListingByMovoto bool   `json:"listingByMovoto"`
	PriceRaw      int    `json:"priceRaw"`
	IsVideoTour   bool   `json:"isVideoTour"`
	Is3DTour      bool   `json:"is3DTour"`
	VideoTourLink string `json:"videoTourLink"`
	VirtualLink   string `json:"virtualLink"`
}

type MovotoJsonRes struct {
	Splits struct {
		MovotoAutopopRtoCW9816                       string `json:"movoto-autopop-rto-CW-9816"`
		MovotoCpclinkCW10060                         string `json:"movoto-cpclink-CW-10060"`
		MovotoEnableDatadogRumCW10249                string `json:"movoto-enable-datadog-rum-CW-10249"`
		MovotoNonRangeLinksCW9935                    string `json:"movoto-non-range-links-CW-9935"`
		MovotoDppPartnerServicesCW10244              string `json:"movoto-dpp-partner-services-CW-10244"`
		MovotoSignupHomeownerCW9793                  string `json:"movoto-signup-homeowner-CW-9793"`
		MovotoMspLessfilterCW10351                   string `json:"movoto-msp-lessfilter-CW-10351"`
		MovotoGeoareaCW10679                         string `json:"movoto-geoarea-CW-10679"`
		MovotoCpcPreapprovedCW10585                  string `json:"movoto-cpc-preapproved-CW-10585"`
		MovotoTextboxCW10373                         string `json:"movoto-textbox-CW-10373"`
		MovotoMobileEngagementCtasStickyOnDppCW10698 string `json:"movoto-mobile-engagement-ctas-sticky-on-dpp-CW-10698"`
		MovotoDppLeadFormCopyCW10590                 string `json:"movoto-dpp-lead-form-copy-CW-10590"`
		MovotoDppFakeChatCW10810                     string `json:"movoto-dpp-fake-chat-CW-10810"`
		MovotoMspCardShareIconCW10805                string `json:"movoto-msp-card-share-icon-CW-10805"`
		MovotoVeteranFunctionalityCW8848             string `json:"movoto-veteran-functionality-CW-8848"`
		MovotoCpcApplyCw9238                         string `json:"movoto-cpc-apply-cw-9238"`
		MovotoCpcCalculateCw9239                     string `json:"movoto-cpc-calculate-cw-9239"`
		MovotoCpcPreapprovedCw9240                   string `json:"movoto-cpc-preapproved-cw-9240"`
		MovotoCpcPrequalifyCw9241                    string `json:"movoto-cpc-prequalify-cw-9241"`
		MovotoCpcLenderCw9242                        string `json:"movoto-cpc-lender-cw-9242"`
		MovotoRentToOwnCw9611                        string `json:"movoto-rent-to-own-cw-9611"`
		MovotoBankrateVeteranTopCW9809               string `json:"movoto-bankrate-veteran-top-CW-9809"`
		MovotoVuVeteranTopCW9812                     string `json:"movoto-vu-veteran-top-CW-9812"`
		MovotoBankrateVeteranBottomCW9810            string `json:"movoto-bankrate-veteran-bottom-CW-9810"`
		MovotoVuVeteranBottomCW9813                  string `json:"movoto-vu-veteran-bottom-CW-9813"`
		MovotoOpendoorSellCW10059                    string `json:"movoto-opendoor-sell-CW-10059"`
		MovotoOpendoorBuyCW10057                     string `json:"movoto-opendoor-buy-CW-10057"`
		MovotoMortgageCW10250                        string `json:"movoto-mortgage-CW-10250"`
		MovotoVeteranTYPCTACW10150                   string `json:"movoto-veteranTYP-CTA-CW-10150"`
		MovotoPostrentalUpdatesCW10243               string `json:"movoto-postrental-updates-CW-10243"`
		MovotoVeteranCtaSplitChangesOnTYPCW10363     string `json:"movoto-veteran-cta-split-changes-on-TYP-CW-10363"`
		MovotoMortgagelongformCW10444                string `json:"movoto-mortgagelongform-CW-10444"`
		MovotoScrollCW10346                          string `json:"movoto-scroll-CW-10346"`
		MovotoIntentfulShortformCw10737              string `json:"movoto-intentful-shortform-cw-10737"`
	} `json:"splits"`
	IP       string `json:"ip"`
	PageData struct {
		RealListingID string      `json:"realListingId"`
		IsSpanishURL  bool        `json:"isSpanishUrl"`
		IsPrOnly      bool        `json:"isPrOnly"`
		ClosePrice    interface{} `json:"closePrice"`
		DaysOnMovoto  int         `json:"daysOnMovoto"`
		Description   string      `json:"description"`
		Features      []struct {
			Name  string `json:"name"`
			Value []struct {
				Name  string `json:"name"`
				Value []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"value"`
			} `json:"value"`
			Icon string `json:"icon,omitempty"`
		} `json:"features"`
		ID           string      `json:"id"`
		TnImgPath    string      `json:"tnImgPath"`
		ListDate     interface{} `json:"listDate"`
		ListingAgent string      `json:"listingAgent"`
		ListPrice    int         `json:"listPrice"`
		LotSize      int         `json:"lotSize"`
		SqftTotal    int         `json:"sqftTotal"`
		MlsDbNumber  int         `json:"mlsDbNumber"`
		Mls          struct {
			CreatedAt                   string      `json:"createdAt"`
			DateHidden                  interface{} `json:"dateHidden"`
			Disclaimer                  string      `json:"disclaimer"`
			DisclaimerPopupText         interface{} `json:"disclaimerPopupText"`
			ExtendedHouseHistory        interface{} `json:"extendedHouseHistory"`
			GeoBoundaryID               int         `json:"geoBoundaryId"`
			ID                          int         `json:"id"`
			LogoURL                     interface{} `json:"logoUrl"`
			MlsName                     string      `json:"mlsName"`
			ShortDisclaimer             string      `json:"shortDisclaimer"`
			ShowAVM                     bool        `json:"showAVM"`
			ShowMovotoAgentSection      bool        `json:"showMovotoAgentSection"`
			ShowPriceHistExtendedFlag   string      `json:"showPriceHistExtendedFlag"`
			ShowPriceHistFlag           string      `json:"showPriceHistFlag"`
			ShowPrivateNotesFlag        bool        `json:"showPrivateNotesFlag"`
			ShowPublicRecordSoldHistory bool        `json:"showPublicRecordSoldHistory"`
			ShowSoldDppFlag             string      `json:"showSoldDppFlag"`
			Type                        string      `json:"type"`
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
		UpdatedTime             string        `json:"updatedTime"`
		HiddenByComplianceRule  bool          `json:"hiddenByComplianceRule"`
		DateHidden              interface{}   `json:"dateHidden"`
		PropertyID              string        `json:"propertyId"`
		Visibility              string        `json:"visibility"`
		PermitAvm               bool          `json:"permitAvm"`
		SoldDate                interface{}   `json:"soldDate"`
		CreatedAt               string        `json:"createdAt"`
		PropertyDateHidden      interface{}   `json:"propertyDateHidden"`
		ImageDownloaderStatus   int           `json:"imageDownloaderStatus"`
		OnMarketDateTime        string        `json:"onMarketDateTime"`
		Garage                  int           `json:"garage"`
		VirtualTourLink         string        `json:"virtualTourLink"`
		VirtualTours            []struct {
			URL                      string `json:"url"`
			VirtualTourDimensionType string `json:"virtualTourDimensionType"`
			VirtualTourType          string `json:"virtualTourType"`
		} `json:"virtualTours"`
		NhsDetails    interface{} `json:"nhsDetails"`
		RentalDetails interface{} `json:"rentalDetails"`
		BuildingName  interface{} `json:"buildingName"`
		AgentDetails  struct {
			AgentContactInfo  interface{} `json:"agentContactInfo"`
			AgentLicenseNo    interface{} `json:"agentLicenseNo"`
			AgentListFullName string      `json:"agentListFullName"`
			AgentListOffice   interface{} `json:"agentListOffice"`
		} `json:"agentDetails"`
		PropertySizeSort int         `json:"propertySizeSort"`
		BrokerageDetails interface{} `json:"brokerageDetails"`
		Geo              struct {
			State              string      `json:"state"`
			City               string      `json:"city"`
			CityID             int         `json:"cityId"`
			County             string      `json:"county"`
			CountyID           int         `json:"countyId"`
			Lat                float64     `json:"lat"`
			Lng                float64     `json:"lng"`
			Zipcode            string      `json:"zipcode"`
			SubPremise         string      `json:"subPremise"`
			Address            string      `json:"address"`
			NeighborhoodNGeoID interface{} `json:"neighborhoodNGeoId"`
			NeighborhoodName   interface{} `json:"neighborhoodName"`
		} `json:"geo"`
		IsNHS       bool   `json:"isNHS"`
		IsRentals   bool   `json:"isRentals"`
		VirtualLink string `json:"virtualLink"`
		FCooling    string `json:"fCooling"`
		FLevels     string `json:"fLevels"`
		IsSold      bool   `json:"isSold"`
		//ListingByMovoto   bool   `json:"listingByMovoto"`
		PriceRaw          int    `json:"priceRaw"`
		IsVideoTour       bool   `json:"isVideoTour"`
		Is3DTour          bool   `json:"is3DTour"`
		VideoTourLink     string `json:"videoTourLink"`
		PublicRecordID    string `json:"publicRecordId"`
		Redirect          bool   `json:"redirect"`
		IsRealPrOnly      bool   `json:"isRealPrOnly"`
		CategorizedPhotos []struct {
			Photos []struct {
				Index  int    `json:"index"`
				URL    string `json:"url"`
				Images struct {
					P string `json:"p"`
					L string `json:"l"`
					R string `json:"r"`
				} `json:"images"`
				MatchRate   int `json:"matchRate"`
				SequenceNum int `json:"sequenceNum"`
			} `json:"photos"`
			Index  int    `json:"index"`
			URL    string `json:"url"`
			Tag    string `json:"tag"`
			Images struct {
				P string `json:"p"`
				L string `json:"l"`
				R string `json:"r"`
			} `json:"images"`
			SequenceNum          int `json:"sequenceNum"`
			MatchedCategoryOrder int `json:"matchedCategoryOrder"`
		} `json:"categorizedPhotos"`
		HasCategories bool `json:"hasCategories"`
		NearbyHomes   []struct {
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
			PriceChange             interface{}   `json:"priceChange"`
			PropertyID              string        `json:"propertyId"`
			Visibility              string        `json:"visibility"`
			SoldDate                interface{}   `json:"soldDate"`
			CreatedAt               string        `json:"createdAt"`
			ImageDownloaderStatus   int           `json:"imageDownloaderStatus"`
			OnMarketDateTime        string        `json:"onMarketDateTime"`
			VirtualTourLink         string        `json:"virtualTourLink"`
			NhsDetails              interface{}   `json:"nhsDetails,omitempty"`
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
			IsNHS     bool `json:"isNHS"`
			IsRentals bool `json:"isRentalhoms"`
			IsSold    bool `json:"isSold"`
			//ListingByMovoto  bool        `json:"listingByMovoto"`
			PriceRaw         int         `json:"priceRaw"`
			IsVideoTour      bool        `json:"isVideoTour"`
			Is3DTour         bool        `json:"is3DTour"`
			VideoTourLink    string      `json:"videoTourLink"`
			VirtualLink      string      `json:"virtualLink,omitempty"`
			Distance         string      `json:"distance"`
			Address          string      `json:"address,omitempty"`
			NhsBuilderName   string      `json:"nhsBuilderName,omitempty"`
			NhsBuilderID     int         `json:"nhsBuilderId,omitempty"`
			NhsListingID     string      `json:"nhsListingId,omitempty"`
			NhsListingType   string      `json:"nhsListingType,omitempty"`
			NhsMarketName    interface{} `json:"nhsMarketName,omitempty"`
			NhsMarketID      interface{} `json:"nhsMarketId,omitempty"`
			NhsPlanName      string      `json:"nhsPlanName,omitempty"`
			NhsCommunityName string      `json:"nhsCommunityName,omitempty"`
			NhsCommunityID   int         `json:"nhsCommunityId,omitempty"`
		} `json:"nearbyHomes"`
		ListingCountObj struct {
			City struct {
				CondoCount                int `json:"condoCount"`
				HomeForSaleCount          int `json:"homeForSaleCount"`
				HomeNewCount              int `json:"homeNewCount"`
				HomeOpenCount             int `json:"homeOpenCount"`
				NewConstructionHomesCount int `json:"newConstructionHomesCount"`
				PoolCount                 int `json:"poolCount"`
				PriceReducedCount         int `json:"priceReducedCount"`
				SchoolCount               int `json:"schoolCount"`
				SingleFamilyCount         int `json:"singleFamilyCount"`
				StateListingCount         int `json:"stateListingCount"`
			} `json:"city"`
			Zipcode struct {
				CondoCount                int `json:"condoCount"`
				HomeForSaleCount          int `json:"homeForSaleCount"`
				HomeNewCount              int `json:"homeNewCount"`
				HomeOpenCount             int `json:"homeOpenCount"`
				NewConstructionHomesCount int `json:"newConstructionHomesCount"`
				PoolCount                 int `json:"poolCount"`
				PriceReducedCount         int `json:"priceReducedCount"`
				SchoolCount               int `json:"schoolCount"`
				SingleFamilyCount         int `json:"singleFamilyCount"`
				StateListingCount         int `json:"stateListingCount"`
			} `json:"zipcode"`
		} `json:"listingCountObj"`
		MarketTrendTable struct {
			City struct {
				AllTypeMedianListPrice                                      int    `json:"allTypeMedianListPrice"`
				AllTypeMedianListPricePercentageCompareWith1YearAgo         string `json:"allTypeMedianListPricePercentageCompareWith1YearAgo"`
				AllTypeMedianPricePerSqftHouse                              int    `json:"allTypeMedianPricePerSqftHouse"`
				AllTypeMedianPricePerSqftHousePercentageCompareWith1YearAgo string `json:"allTypeMedianPricePerSqftHousePercentageCompareWith1YearAgo"`
				AllTypeTotalInventory                                       int    `json:"allTypeTotalInventory"`
				AllTypeTotalInventoryPercentageCompareWith1YearAgo          string `json:"allTypeTotalInventoryPercentageCompareWith1YearAgo"`
				AllTypeMedianDom                                            int    `json:"allTypeMedianDom"`
				AllTypeMedianDomPercentageCompareWith1YearAgo               string `json:"allTypeMedianDomPercentageCompareWith1YearAgo"`
			} `json:"city"`
			Zipcode struct {
				AllTypeMedianListPrice                                      int    `json:"allTypeMedianListPrice"`
				AllTypeMedianListPricePercentageCompareWith1YearAgo         string `json:"allTypeMedianListPricePercentageCompareWith1YearAgo"`
				AllTypeMedianPricePerSqftHouse                              int    `json:"allTypeMedianPricePerSqftHouse"`
				AllTypeMedianPricePerSqftHousePercentageCompareWith1YearAgo string `json:"allTypeMedianPricePerSqftHousePercentageCompareWith1YearAgo"`
				AllTypeTotalInventory                                       int    `json:"allTypeTotalInventory"`
				AllTypeTotalInventoryPercentageCompareWith1YearAgo          string `json:"allTypeTotalInventoryPercentageCompareWith1YearAgo"`
				AllTypeMedianDom                                            int    `json:"allTypeMedianDom"`
				AllTypeMedianDomPercentageCompareWith1YearAgo               string `json:"allTypeMedianDomPercentageCompareWith1YearAgo"`
			} `json:"zipcode"`
			ZipcodeOnemonthago struct {
				AllTypeMedianListPrice         int `json:"allTypeMedianListPrice"`
				AllTypeMedianPricePerSqftHouse int `json:"allTypeMedianPricePerSqftHouse"`
				AllTypeTotalInventory          int `json:"allTypeTotalInventory"`
				AllTypeMedianDom               int `json:"allTypeMedianDom"`
			} `json:"zipcode_onemonthago"`
			Neighborhood interface{} `json:"neighborhood"`
			County       struct {
				AllTypeMedianListPrice         int `json:"allTypeMedianListPrice"`
				AllTypeMedianPricePerSqftHouse int `json:"allTypeMedianPricePerSqftHouse"`
				AllTypeTotalInventory          int `json:"allTypeTotalInventory"`
				AllTypeMedianDom               int `json:"allTypeMedianDom"`
			} `json:"county"`
		} `json:"marketTrendTable"`
		CityMarketSnapshot struct {
			Today struct {
				AllTypeDistressedPercent                                        string `json:"allTypeDistressedPercent"`
				AllTypeMedianDom                                                int    `json:"allTypeMedianDom"`
				AllTypeMedianDomPercentageCompareWith1MonthAgo                  string `json:"allTypeMedianDomPercentageCompareWith1MonthAgo"`
				AllTypeMedianDomPercentageCompareWith1YearAgo                   string `json:"allTypeMedianDomPercentageCompareWith1YearAgo"`
				AllTypeMedianHouseSize                                          int    `json:"allTypeMedianHouseSize"`
				AllTypeMedianHouseSizePercentageCompareWith1MonthAgo            string `json:"allTypeMedianHouseSizePercentageCompareWith1MonthAgo"`
				AllTypeMedianHouseSizePercentageCompareWith1YearAgo             string `json:"allTypeMedianHouseSizePercentageCompareWith1YearAgo"`
				AllTypeMedianListPrice                                          int    `json:"allTypeMedianListPrice"`
				AllTypeMedianListPricePercentageCompareWith1MonthAgo            string `json:"allTypeMedianListPricePercentageCompareWith1MonthAgo"`
				AllTypeMedianListPricePercentageCompareWith1YearAgo             string `json:"allTypeMedianListPricePercentageCompareWith1YearAgo"`
				AllTypeMedianPricePerSqftHouse                                  int    `json:"allTypeMedianPricePerSqftHouse"`
				AllTypeMedianPricePerSqftHousePercentageCompareWith1MonthAgo    string `json:"allTypeMedianPricePerSqftHousePercentageCompareWith1MonthAgo"`
				AllTypeMedianPricePerSqftHousePercentageCompareWith1YearAgo     string `json:"allTypeMedianPricePerSqftHousePercentageCompareWith1YearAgo"`
				AllTypeTotalInventory                                           int    `json:"allTypeTotalInventory"`
				AllTypeTotalInventoryPercentageCompareWith1MonthAgo             string `json:"allTypeTotalInventoryPercentageCompareWith1MonthAgo"`
				AllTypeTotalInventoryPercentageCompareWith1YearAgo              string `json:"allTypeTotalInventoryPercentageCompareWith1YearAgo"`
				CondoTypeDistressedPercent                                      string `json:"condoTypeDistressedPercent"`
				CondoTypeMedianDom                                              int    `json:"condoTypeMedianDom"`
				CondoTypeMedianDomPercentageCompareWith1MonthAgo                string `json:"condoTypeMedianDomPercentageCompareWith1MonthAgo"`
				CondoTypeMedianDomPercentageCompareWith1YearAgo                 string `json:"condoTypeMedianDomPercentageCompareWith1YearAgo"`
				CondoTypeMedianHouseSize                                        int    `json:"condoTypeMedianHouseSize"`
				CondoTypeMedianHouseSizePercentageCompareWith1MonthAgo          string `json:"condoTypeMedianHouseSizePercentageCompareWith1MonthAgo"`
				CondoTypeMedianHouseSizePercentageCompareWith1YearAgo           string `json:"condoTypeMedianHouseSizePercentageCompareWith1YearAgo"`
				CondoTypeMedianListPrice                                        int    `json:"condoTypeMedianListPrice"`
				CondoTypeMedianListPricePercentageCompareWith1MonthAgo          string `json:"condoTypeMedianListPricePercentageCompareWith1MonthAgo"`
				CondoTypeMedianListPricePercentageCompareWith1YearAgo           string `json:"condoTypeMedianListPricePercentageCompareWith1YearAgo"`
				CondoTypeMedianPricePerSqftHouse                                int    `json:"condoTypeMedianPricePerSqftHouse"`
				CondoTypeMedianPricePerSqftHousePercentageCompareWith1MonthAgo  string `json:"condoTypeMedianPricePerSqftHousePercentageCompareWith1MonthAgo"`
				CondoTypeMedianPricePerSqftHousePercentageCompareWith1YearAgo   string `json:"condoTypeMedianPricePerSqftHousePercentageCompareWith1YearAgo"`
				CondoTypeTotalInventory                                         int    `json:"condoTypeTotalInventory"`
				CondoTypeTotalInventoryPercentageCompareWith1MonthAgo           string `json:"condoTypeTotalInventoryPercentageCompareWith1MonthAgo"`
				CondoTypeTotalInventoryPercentageCompareWith1YearAgo            string `json:"condoTypeTotalInventoryPercentageCompareWith1YearAgo"`
				FormatedDate                                                    string `json:"formatedDate"`
				HighestSqftHouse                                                int    `json:"highestSqftHouse"`
				LeastExpensiveHouse                                             int    `json:"leastExpensiveHouse"`
				LowestSqftHouse                                                 int    `json:"lowestSqftHouse"`
				MostExpensiveHouse                                              int    `json:"mostExpensiveHouse"`
				SecondExpensiveHouse                                            int    `json:"secondExpensiveHouse"`
				SingleTypeDistressedPercent                                     string `json:"singleTypeDistressedPercent"`
				SingleTypeMedianDom                                             int    `json:"singleTypeMedianDom"`
				SingleTypeMedianDomPercentageCompareWith1MonthAgo               string `json:"singleTypeMedianDomPercentageCompareWith1MonthAgo"`
				SingleTypeMedianDomPercentageCompareWith1YearAgo                string `json:"singleTypeMedianDomPercentageCompareWith1YearAgo"`
				SingleTypeMedianHouseSize                                       int    `json:"singleTypeMedianHouseSize"`
				SingleTypeMedianHouseSizePercentageCompareWith1MonthAgo         string `json:"singleTypeMedianHouseSizePercentageCompareWith1MonthAgo"`
				SingleTypeMedianHouseSizePercentageCompareWith1YearAgo          string `json:"singleTypeMedianHouseSizePercentageCompareWith1YearAgo"`
				SingleTypeMedianListPrice                                       int    `json:"singleTypeMedianListPrice"`
				SingleTypeMedianListPricePercentageCompareWith1MonthAgo         string `json:"singleTypeMedianListPricePercentageCompareWith1MonthAgo"`
				SingleTypeMedianListPricePercentageCompareWith1YearAgo          string `json:"singleTypeMedianListPricePercentageCompareWith1YearAgo"`
				SingleTypeMedianPricePerSqftHouse                               int    `json:"singleTypeMedianPricePerSqftHouse"`
				SingleTypeMedianPricePerSqftHousePercentageCompareWith1MonthAgo string `json:"singleTypeMedianPricePerSqftHousePercentageCompareWith1MonthAgo"`
				SingleTypeMedianPricePerSqftHousePercentageCompareWith1YearAgo  string `json:"singleTypeMedianPricePerSqftHousePercentageCompareWith1YearAgo"`
				SingleTypeTotalInventory                                        int    `json:"singleTypeTotalInventory"`
				SingleTypeTotalInventoryPercentageCompareWith1MonthAgo          string `json:"singleTypeTotalInventoryPercentageCompareWith1MonthAgo"`
				SingleTypeTotalInventoryPercentageCompareWith1YearAgo           string `json:"singleTypeTotalInventoryPercentageCompareWith1YearAgo"`
			} `json:"today"`
			OneMonthAgo struct {
				AllTypeDistressedPercent          string `json:"allTypeDistressedPercent"`
				AllTypeMedianDom                  int    `json:"allTypeMedianDom"`
				AllTypeMedianHouseSize            int    `json:"allTypeMedianHouseSize"`
				AllTypeMedianListPrice            int    `json:"allTypeMedianListPrice"`
				AllTypeMedianPricePerSqftHouse    int    `json:"allTypeMedianPricePerSqftHouse"`
				AllTypeTotalInventory             int    `json:"allTypeTotalInventory"`
				CondoTypeDistressedPercent        string `json:"condoTypeDistressedPercent"`
				CondoTypeMedianDom                int    `json:"condoTypeMedianDom"`
				CondoTypeMedianHouseSize          int    `json:"condoTypeMedianHouseSize"`
				CondoTypeMedianListPrice          int    `json:"condoTypeMedianListPrice"`
				CondoTypeMedianPricePerSqftHouse  int    `json:"condoTypeMedianPricePerSqftHouse"`
				CondoTypeTotalInventory           int    `json:"condoTypeTotalInventory"`
				FormatedDate                      string `json:"formatedDate"`
				HighestSqftHouse                  int    `json:"highestSqftHouse"`
				LeastExpensiveHouse               int    `json:"leastExpensiveHouse"`
				LowestSqftHouse                   int    `json:"lowestSqftHouse"`
				MostExpensiveHouse                int    `json:"mostExpensiveHouse"`
				SecondExpensiveHouse              int    `json:"secondExpensiveHouse"`
				SingleTypeDistressedPercent       string `json:"singleTypeDistressedPercent"`
				SingleTypeMedianDom               int    `json:"singleTypeMedianDom"`
				SingleTypeMedianHouseSize         int    `json:"singleTypeMedianHouseSize"`
				SingleTypeMedianListPrice         int    `json:"singleTypeMedianListPrice"`
				SingleTypeMedianPricePerSqftHouse int    `json:"singleTypeMedianPricePerSqftHouse"`
				SingleTypeTotalInventory          int    `json:"singleTypeTotalInventory"`
			} `json:"one_month_ago"`
			OneYearAgo struct {
				AllTypeDistressedPercent          string `json:"allTypeDistressedPercent"`
				AllTypeMedianDom                  int    `json:"allTypeMedianDom"`
				AllTypeMedianHouseSize            int    `json:"allTypeMedianHouseSize"`
				AllTypeMedianListPrice            int    `json:"allTypeMedianListPrice"`
				AllTypeMedianPricePerSqftHouse    int    `json:"allTypeMedianPricePerSqftHouse"`
				AllTypeTotalInventory             int    `json:"allTypeTotalInventory"`
				CondoTypeDistressedPercent        string `json:"condoTypeDistressedPercent"`
				CondoTypeMedianDom                int    `json:"condoTypeMedianDom"`
				CondoTypeMedianHouseSize          int    `json:"condoTypeMedianHouseSize"`
				CondoTypeMedianListPrice          int    `json:"condoTypeMedianListPrice"`
				CondoTypeMedianPricePerSqftHouse  int    `json:"condoTypeMedianPricePerSqftHouse"`
				CondoTypeTotalInventory           int    `json:"condoTypeTotalInventory"`
				FormatedDate                      string `json:"formatedDate"`
				HighestSqftHouse                  int    `json:"highestSqftHouse"`
				LeastExpensiveHouse               int    `json:"leastExpensiveHouse"`
				LowestSqftHouse                   int    `json:"lowestSqftHouse"`
				MostExpensiveHouse                int    `json:"mostExpensiveHouse"`
				SecondExpensiveHouse              int    `json:"secondExpensiveHouse"`
				SingleTypeDistressedPercent       string `json:"singleTypeDistressedPercent"`
				SingleTypeMedianDom               int    `json:"singleTypeMedianDom"`
				SingleTypeMedianHouseSize         int    `json:"singleTypeMedianHouseSize"`
				SingleTypeMedianListPrice         int    `json:"singleTypeMedianListPrice"`
				SingleTypeMedianPricePerSqftHouse int    `json:"singleTypeMedianPricePerSqftHouse"`
				SingleTypeTotalInventory          int    `json:"singleTypeTotalInventory"`
			} `json:"one_year_ago"`
		} `json:"cityMarketSnapshot"`
		LocalHighlight struct {
			Sections []struct {
				SectionTitle string `json:"sectionTitle"`
				List         []struct {
					Title          string `json:"title"`
					Value          string `json:"value"`
					Image          string `json:"image"`
					Description    string `json:"description,omitempty"`
					MonthlyAverage struct {
						Low         []int    `json:"low"`
						High        []int    `json:"high"`
						Months      []string `json:"months"`
						WarmestText string   `json:"warmestText"`
						CoolestText string   `json:"coolestText"`
					} `json:"monthlyAverage,omitempty"`
				} `json:"list"`
			} `json:"sections"`
			SubHeader string `json:"subHeader"`
			Property  struct {
				Propertyid string `json:"propertyid"`
				Sqft       int    `json:"sqft"`
				YearBuilt  int    `json:"year_built"`
			} `json:"property"`
			NearestTransit string `json:"nearestTransit"`
		} `json:"localHighlight"`
		LocalHighlightScore struct {
			Foodie struct {
				Score    int `json:"score"`
				Quantity struct {
					Restaurants int `json:"restaurants"`
					Density     int `json:"density"`
				} `json:"quantity"`
			} `json:"foodie"`
			Kids struct {
				Score   int `json:"score"`
				Schools struct {
					Total  int `json:"total"`
					Impact int `json:"impact"`
				} `json:"schools"`
				Parks struct {
					Total  int `json:"total"`
					Impact int `json:"impact"`
				} `json:"parks"`
				Pois struct {
					Total  int `json:"total"`
					Impact int `json:"impact"`
				} `json:"pois"`
			} `json:"kids"`
			Dogs struct {
				Score  int `json:"score"`
				Trails struct {
					Miles  int `json:"miles"`
					Total  int `json:"total"`
					Impact int `json:"impact"`
				} `json:"trails"`
				Parks struct {
					Total  int `json:"total"`
					Impact int `json:"impact"`
				} `json:"parks"`
				Pois struct {
					Total  int `json:"total"`
					Impact int `json:"impact"`
				} `json:"pois"`
			} `json:"dogs"`
			WalkscoreData struct {
				Score     int `json:"score"`
				Walkscore []struct {
					ID     string `json:"id"`
					Label  string `json:"label"`
					Total  int    `json:"total"`
					Impact int    `json:"impact"`
				} `json:"walkscore"`
			} `json:"walkscore_data"`
		} `json:"localHighlightScore"`
		ClimateListData struct {
			List []struct {
				Title       string `json:"title"`
				Level       string `json:"level"`
				Score       int    `json:"score"`
				Description string `json:"description"`
			} `json:"list"`
			AvgClimate  string `json:"avgClimate"`
			Description string `json:"description"`
		} `json:"climateListData"`
		PriceHistory []struct {
			Change         string `json:"change"`
			DataSource     string `json:"dataSource"`
			Date           string `json:"date"`
			HistoryID      string `json:"historyId"`
			MlsID          int    `json:"mlsId"`
			MlsNumber      string `json:"mlsNumber"`
			Price          int    `json:"price"`
			Status         string `json:"status"`
			DateDesktop    string `json:"dateDesktop"`
			DateMobile     string `json:"dateMobile"`
			PriceDesktop   string `json:"priceDesktop"`
			PriceMobile    string `json:"priceMobile"`
			PriceChangeTag int    `json:"priceChangeTag,omitempty"`
			IsRental       bool   `json:"isRental"`
		} `json:"priceHistory"`
		IsOjOBrokerage bool   `json:"isOjOBrokerage"`
		GeoPhone       string `json:"geoPhone"`
		UpdatedStatus  struct {
			ModificationTimestamp string `json:"modificationTimestamp"`
			LastUpdate            string `json:"lastUpdate"`
			LastChecked           string `json:"lastChecked"`
		} `json:"updatedStatus"`
		SeoNearbyCity []struct {
			ID                         int    `json:"id"`
			DisplayName                string `json:"displayName"`
			Type                       string `json:"type"`
			StateCode                  string `json:"stateCode"`
			StateDisplayName           string `json:"stateDisplayName"`
			SitemapCityPropertyPageURL string `json:"sitemapCityPropertyPageUrl"`
			MedianListPrice            int    `json:"medianListPrice"`
			Rank                       int    `json:"rank"`
			Location                   struct {
				Lng float64 `json:"lng"`
				Lat float64 `json:"lat"`
			} `json:"location"`
			MedianListPriceStr string `json:"medianListPriceStr"`
			Distance           string `json:"distance"`
			CarCommute         string `json:"carCommute"`
		} `json:"seoNearbyCity"`
		TopNeighborhoods []struct {
			ID                         int    `json:"id"`
			CityID                     int    `json:"cityId"`
			CityDisplayName            string `json:"cityDisplayName"`
			DisplayName                string `json:"displayName"`
			Type                       string `json:"type"`
			StateCode                  string `json:"stateCode"`
			StateDisplayName           string `json:"stateDisplayName"`
			SitemapCityPropertyPageURL string `json:"sitemapCityPropertyPageUrl"`
			SitemapNeighborhoodPageURL string `json:"sitemapNeighborhoodPageUrl"`
			MedianListPrice            int    `json:"medianListPrice"`
			Rank                       int    `json:"rank"`
			Location                   struct {
				Lng float64 `json:"lng"`
				Lat float64 `json:"lat"`
			} `json:"location"`
			MedianListPriceStr string `json:"medianListPriceStr"`
		} `json:"topNeighborhoods"`
		SeoNearbyZipCode []struct {
			ID                         int    `json:"id"`
			CityID                     int    `json:"cityId"`
			CityDisplayName            string `json:"cityDisplayName"`
			DisplayName                string `json:"displayName"`
			Type                       string `json:"type"`
			StateCode                  string `json:"stateCode"`
			StateDisplayName           string `json:"stateDisplayName"`
			SitemapCityPropertyPageURL string `json:"sitemapCityPropertyPageUrl"`
			SitemapZipCodePageURL      string `json:"sitemapZipCodePageUrl"`
			MedianListPrice            int    `json:"medianListPrice"`
			Rank                       int    `json:"rank"`
			Location                   struct {
				Lng float64 `json:"lng"`
				Lat float64 `json:"lat"`
			} `json:"location"`
			MedianListPriceStr string `json:"medianListPriceStr"`
		} `json:"seoNearbyZipCode"`
		TopNearbyCounty []struct {
			ID                   int    `json:"id"`
			DisplayName          string `json:"displayName"`
			Type                 string `json:"type"`
			StateCode            string `json:"stateCode"`
			StateDisplayName     string `json:"stateDisplayName"`
			SitemapCountyPageURL string `json:"sitemapCountyPageUrl"`
			MedianListPrice      int    `json:"medianListPrice"`
			Rank                 int    `json:"rank"`
			Location             struct {
				Lng float64 `json:"lng"`
				Lat float64 `json:"lat"`
			} `json:"location"`
			MedianListPriceStr string `json:"medianListPriceStr"`
		} `json:"topNearbyCounty"`
		NeighborhoodInfo     interface{} `json:"neighborhoodInfo"`
		ComplianceListOffice string      `json:"complianceListOffice"`
		NearbySchools        []struct {
			Distance            string      `json:"distance"`
			DspURL              string      `json:"dspUrl"`
			ID                  string      `json:"id"`
			Latitude            float64     `json:"latitude"`
			Level               string      `json:"level"`
			LevelCode           string      `json:"levelCode"`
			Longitude           float64     `json:"longitude"`
			NcesID              string      `json:"ncesId"`
			Rating              int         `json:"rating"`
			ReviewCount         int         `json:"reviewCount"`
			Name                string      `json:"name"`
			Type                string      `json:"type"`
			ReveiwAvgQuality    int         `json:"reveiwAvgQuality"`
			ReviewAvgActivities interface{} `json:"reviewAvgActivities"`
			ReviewAvgParents    interface{} `json:"reviewAvgParents"`
			ReviewAvgPrincipal  interface{} `json:"reviewAvgPrincipal"`
			ReviewAvgSafety     interface{} `json:"reviewAvgSafety"`
			ReviewAvgTeachers   interface{} `json:"reviewAvgTeachers"`
		} `json:"nearbySchools"`
		NearbySchoolAvgRating struct {
			Min int `json:"min"`
			Max int `json:"max"`
		} `json:"nearbySchoolAvgRating"`
		SchoolDistricts []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			NcesID      string `json:"ncesId"`
			DistrictURL string `json:"districtUrl"`
		} `json:"schoolDistricts"`
		PublicRecord struct {
			AirConditioningDesc  string      `json:"airConditioningDesc"`
			BasementDesc         string      `json:"basementDesc"`
			BuildingArea         string      `json:"buildingArea"`
			ElevatorDesc         string      `json:"elevatorDesc"`
			EstValue             string      `json:"estValue"`
			ExteriorWallsDesc    string      `json:"exteriorWallsDesc"`
			Fireplace            string      `json:"fireplace"`
			FoundationDesc       string      `json:"foundationDesc"`
			GarageTypeDesc       string      `json:"garageTypeDesc"`
			HeatingDesc          string      `json:"heatingDesc"`
			LandUseCodeDesc      string      `json:"landUseCodeDesc"`
			LastSaleDate         string      `json:"lastSaleDate"`
			LastSaleDateRaw      interface{} `json:"lastSaleDateRaw"`
			LastSalePrice        string      `json:"lastSalePrice"`
			LotSize              string      `json:"lotSize"`
			LotSizeUnit          string      `json:"lotSizeUnit"`
			NumBaths             int         `json:"numBaths"`
			NumBeds              int         `json:"numBeds"`
			NumGarages           int         `json:"numGarages"`
			NumPartBathsDesc     string      `json:"numPartBathsDesc"`
			NumRooms             string      `json:"numRooms"`
			PoolDesc             string      `json:"poolDesc"`
			RoofCoverDesc        string      `json:"roofCoverDesc"`
			StyleDesc            string      `json:"styleDesc"`
			TypeConstructionDesc string      `json:"typeConstructionDesc"`
			YearBuilt            int         `json:"yearBuilt"`
			BuildingAreaUnit     string      `json:"buildingAreaUnit"`
		} `json:"publicRecord"`
		Faqs []struct {
			Question string `json:"question"`
			Answer   string `json:"answer"`
		} `json:"faqs"`
		PropertyAttributes []struct {
			AttributeID string `json:"attributeId"`
			DisplayName string `json:"displayName"`
			Rank        int    `json:"rank"`
		} `json:"propertyAttributes"`
		IsHot           bool `json:"isHot"`
		HotProbability  int  `json:"hotProbability"`
		IsCheapProperty bool `json:"isCheapProperty"`
		NearbyEstPrice  struct {
			EstPrice          int     `json:"estPrice"`
			Min               int     `json:"min"`
			Max               int     `json:"max"`
			AreaPriceAvg      int     `json:"areaPriceAvg"`
			Area              int     `json:"area"`
			Percentage        float64 `json:"percentage"`
			IsAbove           bool    `json:"isAbove"`
			EstPriceRangeText string  `json:"estPriceRangeText"`
		} `json:"nearbyEstPrice"`
		Head struct {
		} `json:"head"`
		Seo struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			H1          string `json:"h1"`
		} `json:"seo"`
		AvgPendingCount int `json:"avgPendingCount"`
	} `json:"pageData"`
	Date                  time.Time `json:"date"`
	PageType              string    `json:"pageType"`
	Version               string    `json:"version"`
	AppURL                string    `json:"appUrl"`
	OjoURL                string    `json:"ojoUrl"`
	CdnURL                string    `json:"cdnUrl"`
	DigsURL               string    `json:"digsUrl"`
	StaticURL             string    `json:"staticUrl"`
	CdnIconURL            string    `json:"cdnIconUrl"`
	NovaAgentDesktopURL   string    `json:"novaAgentDesktopUrl"`
	SignalAPIURL          string    `json:"signalAPIUrl"`
	BoundaryURL           string    `json:"boundaryUrl"`
	RequestURLRaw         string    `json:"requestUrlRaw"`
	ReferURL              string    `json:"referUrl"`
	URL                   string    `json:"url"`
	Flat                  int       `json:"flat"`
	PhoneNumber           string    `json:"phoneNumber"`
	GoogleMapKey          string    `json:"googleMapKey"`
	GaAccount             string    `json:"gaAccount"`
	Ga4Measurement        string    `json:"ga4Measurement"`
	GtmContainer          string    `json:"gtmContainer"`
	MarkUserViewedTimeout int       `json:"markUserViewedTimeout"`
	FirebaseMessaging     struct {
		Chrome struct {
			APIKey            string `json:"apiKey"`
			AppID             string `json:"appId"`
			AuthDomain        string `json:"authDomain"`
			DatabaseURL       string `json:"databaseURL"`
			ProjectID         string `json:"projectId"`
			StorageBucket     string `json:"storageBucket"`
			MessagingSenderID string `json:"messagingSenderId"`
			PublicKey         string `json:"publicKey"`
			Debug             bool   `json:"debug"`
		} `json:"chrome"`
		Safari struct {
			PushID  string `json:"pushId"`
			Package string `json:"package"`
		} `json:"safari"`
	} `json:"firebaseMessaging"`
	GoogleLoginClientID string `json:"googleLoginClientId"`
	AppleLoginClientID  string `json:"appleLoginClientId"`
	IsDevelopment       bool   `json:"isDevelopment"`
	IsUserBrowser       bool   `json:"isUserBrowser"`
	Brokerages          struct {
		Cities []struct {
			ClassName   string `json:"className"`
			DisplayName string `json:"displayName"`
			StateCode   string `json:"stateCode"`
			CityID      string `json:"cityId"`
		} `json:"cities"`
		Office struct {
			CA struct {
				Num67121 struct {
					DisplayName  string `json:"displayName"`
					Address      string `json:"address"`
					Phone        string `json:"phone"`
					GoogleMapURL string `json:"googleMapUrl"`
					Location     struct {
						Lat string `json:"lat"`
						Lng string `json:"lng"`
					} `json:"location"`
					Image struct {
						Xs  string `json:"xs"`
						Sm  string `json:"sm"`
						Alt string `json:"alt"`
					} `json:"image"`
					URL string `json:"url"`
				} `json:"67121"`
				Num69054 struct {
					DisplayName  string `json:"displayName"`
					Address      string `json:"address"`
					Phone        string `json:"phone"`
					GoogleMapURL string `json:"googleMapUrl"`
					Location     struct {
						Lat string `json:"lat"`
						Lng string `json:"lng"`
					} `json:"location"`
					Image struct {
						Xs  string `json:"xs"`
						Sm  string `json:"sm"`
						Alt string `json:"alt"`
					} `json:"image"`
					URL string `json:"url"`
				} `json:"69054"`
				Num69270 struct {
					DisplayName  string `json:"displayName"`
					Address      string `json:"address"`
					Phone        string `json:"phone"`
					GoogleMapURL string `json:"googleMapUrl"`
					Location     struct {
						Lat string `json:"lat"`
						Lng string `json:"lng"`
					} `json:"location"`
					Image struct {
						Xs  string `json:"xs"`
						Sm  string `json:"sm"`
						Alt string `json:"alt"`
					} `json:"image"`
					URL string `json:"url"`
				} `json:"69270"`
				Num69324 struct {
					DisplayName  string `json:"displayName"`
					Phone        string `json:"phone"`
					GoogleMapURL string `json:"googleMapUrl"`
					Location     struct {
						Lat string `json:"lat"`
						Lng string `json:"lng"`
					} `json:"location"`
					Image struct {
						Xs  string `json:"xs"`
						Sm  string `json:"sm"`
						Alt string `json:"alt"`
					} `json:"image"`
					URL string `json:"url"`
				} `json:"69324"`
				Num69797 struct {
					DisplayName  string `json:"displayName"`
					Address      string `json:"address"`
					Phone        string `json:"phone"`
					GoogleMapURL string `json:"googleMapUrl"`
					Location     struct {
						Lat string `json:"lat"`
						Lng string `json:"lng"`
					} `json:"location"`
					Image struct {
						Xs  string `json:"xs"`
						Sm  string `json:"sm"`
						Alt string `json:"alt"`
					} `json:"image"`
					URL string `json:"url"`
				} `json:"69797"`
			} `json:"CA"`
			NV struct {
				Num66849 struct {
					DisplayName  string `json:"displayName"`
					Address      string `json:"address"`
					Phone        string `json:"phone"`
					GoogleMapURL string `json:"googleMapUrl"`
					Location     struct {
						Lat string `json:"lat"`
						Lng string `json:"lng"`
					} `json:"location"`
					Image struct {
						Xs  string `json:"xs"`
						Sm  string `json:"sm"`
						Alt string `json:"alt"`
					} `json:"image"`
					URL string `json:"url"`
				} `json:"66849"`
			} `json:"NV"`
			AZ struct {
				Num65543 struct {
					DisplayName  string `json:"displayName"`
					Address      string `json:"address"`
					Phone        string `json:"phone"`
					GoogleMapURL string `json:"googleMapUrl"`
					Location     struct {
						Lat string `json:"lat"`
						Lng string `json:"lng"`
					} `json:"location"`
					Image struct {
						Xs  string `json:"xs"`
						Sm  string `json:"sm"`
						Alt string `json:"alt"`
					} `json:"image"`
					URL string `json:"url"`
				} `json:"65543"`
			} `json:"AZ"`
		} `json:"office"`
	} `json:"brokerages"`
	IsPhone        bool   `json:"isPhone"`
	IsMobile       bool   `json:"isMobile"`
	Language       string `json:"language"`
	IsSafari       bool   `json:"isSafari"`
	OS             string `json:"OS"`
	EnableVWO      bool   `json:"enableVWO"`
	RandomID       int    `json:"randomId"`
	MovotoDeviceID string `json:"movotoDeviceID"`
	InternalUser   bool   `json:"internalUser"`
	MlsLogoIds     []int  `json:"mlsLogoIds"`
	Runat          string `json:"runat"`
	Seo            struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		H1          string `json:"h1"`
	} `json:"seo"`
	User struct {
		ID        string `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Role      string `json:"role"`
		OjoID     string `json:"ojoId"`
	} `json:"user"`
	AssignedAgents []interface{} `json:"assignedAgents"`
	Agent          struct {
		ID                    string `json:"id"`
		Name                  string `json:"name"`
		Type                  string `json:"type"`
		Email                 string `json:"email"`
		Phone                 string `json:"phone"`
		IsAgentFeedback       string `json:"isAgentFeedback"`
		RelationshipCreatedAt string `json:"relationshipCreatedAt"`
	} `json:"agent"`
	HotleadInfo          interface{} `json:"hotleadInfo"`
	EnableThirdPart      bool        `json:"enableThirdPart"`
	EnableSW             bool        `json:"enableSW"`
	EnablePerformanceLog bool        `json:"enablePerformanceLog"`
	AbConf               struct {
		CW6284 []string `json:"CW-6284"`
	} `json:"abConf"`
	Simple         bool     `json:"simple"`
	Sematic        bool     `json:"sematic"`
	ExtendedOS     string   `json:"extendedOS"`
	PreloadGa      bool     `json:"preloadGa"`
	Rentals        bool     `json:"rentals"`
	RentLang       string   `json:"rentLang"`
	SplitWhiteList []string `json:"splitWhiteList"`
	SignalJWT      string   `json:"signalJWT"`
	PageInfo       struct {
		UserType        string `json:"userType"`
		FullPageType    string `json:"fullPageType"`
		Value           string `json:"value"`
		Sitesection     string `json:"sitesection"`
		Type            string `json:"type"`
		ListingPagetype string `json:"listing_pagetype"`
	} `json:"pageInfo"`
	SemFlag bool `json:"semFlag"`
	Geo     struct {
		State              string      `json:"state"`
		City               string      `json:"city"`
		CityID             int         `json:"cityId"`
		County             string      `json:"county"`
		CountyID           int         `json:"countyId"`
		Lat                float64     `json:"lat"`
		Lng                float64     `json:"lng"`
		Zipcode            string      `json:"zipcode"`
		SubPremise         string      `json:"subPremise"`
		Address            string      `json:"address"`
		NeighborhoodNGeoID interface{} `json:"neighborhoodNGeoId"`
		NeighborhoodName   interface{} `json:"neighborhoodName"`
	} `json:"geo"`
	IsWebViewMode    bool     `json:"isWebViewMode"`
	TransactionCount int      `json:"transactionCount"`
	AgentCount       int      `json:"agentCount"`
	WelcomeImage     string   `json:"welcomeImage"`
	SplitClasses     []string `json:"splitClasses"`
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
	HOAFee                     int
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
	LotSizeSF                  int
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
