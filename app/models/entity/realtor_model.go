package entity

import (
	"gorm.io/gorm"
	"multilogin_scraping/app/models/base"
	util "multilogin_scraping/pkg/utils"
	"time"
)

// TODO: Use swagger later
type Realtor struct {
	base.BaseIDModel
	// Add pointer to foreign key to set null
	//Maindb3ID                  uint64    `gorm:"column:maindb3_id; unique" json:"maindb3_id"`
	URL                        string    `gorm:"type:text" json:"url"`
	Address                    string    `gorm:"type:text" json:"address"`
	PropertyStatus             bool      `json:"property_status"`
	Bed                        float64   `json:"bed"`
	Bath                       float64   `gorm:"type:decimal" json:"bath"`
	FullBathrooms              float64   `gorm:"type:decimal" json:"full_bathrooms"`
	HalfBathrooms              float64   `gorm:"type:decimal" json:"half_bathrooms"`
	SF                         float64   `gorm:"type:decimal" json:"sf"`
	SalesPrice                 float64   `gorm:"type:decimal" json:"sales_price"`
	EstPayment                 string    `gorm:"type:varchar(250)" json:"est_payment"`
	PrincipalInterest          string    `gorm:"type:varchar(250)" json:"principal_interest"`
	MortgageInsurance          string    `gorm:"type:varchar(250)" json:"mortgage_insurance"`
	PropertyTaxes              string    `gorm:"type:varchar(250)" json:"property_taxes"`
	HomeInsurance              string    `gorm:"type:varchar(250)" json:"home_insurance"`
	HOAFee                     string    `gorm:"type:varchar(250)" json:"hoa_fee"`
	Utilities                  string    `gorm:"type:varchar(250)" json:"utilities"`
	RentZestimate              float64   `gorm:"type:decimal" json:"rent_zestimate"`
	Zestimate                  float64   `gorm:"type:decimal" json:"zestimate"`
	EstimatedSalesRangeMinimum string    `gorm:"type:varchar(250)" json:"estimated_sales_range_minimum"`
	EstimatedSalesRangeMax     string    `gorm:"type:varchar(250)" json:"estimated_sales_range_max"`
	Pictures                   string    `json:"pictures"`
	TimeOnZillow               string    `gorm:"type:varchar(250)" json:"time_on_zillow"`
	Views                      int       `json:"views"`
	Saves                      int       `json:"saves"`
	Overview                   string    `gorm:"type:text" json:"overview"`
	ZillowCheckedDate          string    `gorm:"type:varchar(250)" json:"zillow_checked_date"`
	DataUploadedDate           string    `gorm:"type:varchar(250)" json:"data_uploaded_date"`
	ListedBy                   string    `json:"listed_by"`
	Source                     string    `gorm:"type:text" json:"source"`
	MLS                        string    `gorm:"type:varchar(250)" json:"mls"`
	PropertyType               string    `gorm:"type:varchar(250)" json:"property_type"`
	YearBuilt                  string    `gorm:"type:varchar(250)" json:"year_built"`
	NaturalGas                 bool      `json:"natural_gas"`
	CentralAir                 bool      `json:"central_air"`
	OfGarageSpaces             string    `gorm:"type:varchar(250)" json:"of_garage_spaces"`
	HOAAmount                  string    `gorm:"type:varchar(250)" json:"hoa_amount"`
	LotSizeSF                  string    `gorm:"type:varchar(250)" json:"lot_size_sf"`
	LotSizeAcres               string    `gorm:"type:varchar(250)" json:"lot_size_acres"`
	BuyerAgentFee              string    `gorm:"type:varchar(250)" json:"buyer_agent_fee"`
	Appliances                 string    `gorm:"type:varchar(250)" json:"appliances"`
	LivingRoomLevel            string    `gorm:"type:varchar(250)" json:"living_room_level"`
	LivingRoomDimensions       string    `gorm:"type:varchar(250)" json:"living_room_dimensions"`
	InteriorFeatures           string    `gorm:"type:text" json:"interior_features"`
	PrimaryBedroomLevel        string    `gorm:"type:varchar(250)" json:"primary_bedroom_level"`
	PrimaryBedroomDimensions   string    `gorm:"type:varchar(250)" json:"primary_bedroom_dimensions"`
	Basement                   string    `gorm:"type:varchar(250)" json:"basement"`
	TotalInteriorLivableAreaSF string    `gorm:"type:varchar(250)" json:"total_interior_livable_area_sf"`
	OfFireplaces               string    `gorm:"type:varchar(250)" json:"of_fireplaces"`
	FireplaceFeatures          string    `gorm:"type:varchar(250)" json:"fireplace_features"`
	FlooringType               string    `gorm:"type:varchar(250)" json:"flooring_type"`
	HeatingType                string    `gorm:"type:varchar(250)" json:"heating_type"`
	TotalParkingSpaces         string    `gorm:"type:varchar(250)" json:"total_parking_spaces"`
	ParkingFeatures            string    `gorm:"type:varchar(250)" json:"parking_features"`
	LotFeatures                string    `gorm:"type:varchar(250)" json:"lot_features"`
	CoveredSpaces              string    `gorm:"type:varchar(250)" json:"covered_spaces"`
	ParcelNumber               string    `gorm:"type:varchar(250)" json:"parcel_number"`
	LevelsStoriesFloors        string    `gorm:"type:varchar(250)" json:"levels_stories_floors"`
	PatioAndPorchDetails       string    `gorm:"type:varchar(250)" json:"patio_and_porch_details"`
	HomeType                   string    `gorm:"type:varchar(250)" json:"home_type"`
	ProperySubType             string    `gorm:"type:varchar(250)" json:"propery_sub_type"`
	ConstructionMaterials      string    `gorm:"type:varchar(250)" json:"construction_materials"`
	Foundation                 string    `gorm:"type:varchar(250)" json:"foundation"`
	Roof                       string    `gorm:"type:varchar(250)" json:"roof"`
	NewConstruction            string    `gorm:"type:varchar(250)" json:"new_construction"`
	SewerInformation           string    `gorm:"type:varchar(250)" json:"sewer_information"`
	WaterInformation           string    `gorm:"type:varchar(250)" json:"water_information"`
	RegionLocation             string    `gorm:"type:varchar(250)" json:"region_location"`
	Subdivision                string    `gorm:"type:varchar(250)" json:"subdivision"`
	HasHOA                     string    `gorm:"type:varchar(250)" json:"has_hoa"`
	HOAFeeDetail               string    `gorm:"type:varchar(250)" json:"hoa_fee_detail"`
	ServicesIncluded           string    `gorm:"type:varchar(250)" json:"services_included"`
	AssociationName            string    `gorm:"type:varchar(250)" json:"association_name"`
	AssociationPhone           string    `gorm:"type:varchar(250)" json:"association_phone"`
	AnnualTaxAmount            string    `gorm:"type:varchar(250)" json:"annual_tax_amount"`
	ElementarySchool           string    `gorm:"type:varchar(250)" json:"elementary_school"`
	MiddleSchool               string    `gorm:"type:varchar(250)" json:"middle_school"`
	HighSchool                 string    `gorm:"type:varchar(250)" json:"high_school"`
	District                   string    `gorm:"type:varchar(250)" json:"district"`
	DataSource                 string    `gorm:"type:varchar(250)" json:"data_source"`
	CountyTaxAssessorURL       string    `gorm:"type:text" json:"county_tax_assessor_url"`
	TimestampForDataExtraction time.Time `gorm:"type:timestamp" json:"timestamp_for_data_extraction"`
	CrawlingStatus             string    `gorm:"type:varchar(50)" json:"crawling_status"`
}

// TableName overrides
func (Realtor) TableName() string {
	return "realtor"
}

func (base *Realtor) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("TimestampForDataExtraction", time.Now().In(util.Loc))
	return nil
}
