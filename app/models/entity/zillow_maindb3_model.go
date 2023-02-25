package entity

import (
	"database/sql"
	"gorm.io/gorm"
)

type ZillowMaindb3Address struct {
	gorm.Model
	ID                    uint64 `gorm:"primary_key:auto_increment" json:"id"`
	URL                   string `gorm:"type:text; unique" json:"url"`
	AddressStreet         string `gorm:"type:text;column:Address_Street" json:"Address_Street"`
	AddressCity           string `gorm:"type:varchar(250);column:Address_City" json:"Address_City"`
	AddressState          string `gorm:"type:varchar(250);column:Address_State" json:"Address_State"`
	AddressZip            string `gorm:"type:varchar(250);column:Address_Zip" json:"Address_Zip"`
	AddressZip4           string `gorm:"type:varchar(250);column:Address_Zip4" json:"Address_Zip4"`
	CompleteAddress       string `gorm:"type:varchar(250);column:Complete_Address" json:"Complete_Address"`
	LotSize               string `gorm:"type:varchar(250);column:Lot_Size" json:"Lot_Size"`
	LandAcres             string `gorm:"type:varchar(250);column:Land_Acres" json:"Land_Acres"`
	ListingPrice          string `gorm:"type:varchar(250);column:Listing_Price" json:"Listing_Price"`
	Beds                  string `gorm:"type:varchar(250);column:Beds" json:"Beds"`
	Baths                 string `gorm:"type:varchar(250);column:Baths" json:"Baths"`
	PropertyType          string `gorm:"type:varchar(250);column:Property_Type" json:"Property_Type"`
	YearBuilt             string `gorm:"type:varchar(250);column:Year_Built" json:"Year_Built"`
	GarageSpaces          string `gorm:"type:varchar(250);column:Garage_Spaces" json:"Garage_Spaces"`
	County                string `gorm:"type:varchar(250);column:County" json:"County"`
	MlsID                 string `gorm:"type:varchar(250);column:MLS_ID" json:"MLS_ID"`
	MovotoURL             string `gorm:"type:text;column:Movoto_URL" json:"Movoto_URL"`
	School                string `gorm:"type:varchar(250);column:School" json:"School"`
	LivingArea            string `gorm:"type:varchar(250);column:Living_Area" json:"Living_Area"`
	Pool                  string `gorm:"type:varchar(1);column:Pool" json:"Pool"`
	CentralHeat           string `gorm:"type:char(1);column:Central_Heat" json:"Central_Heat"`
	CentralAir            string `gorm:"type:char(1);column:Central_Air" json:"Central_Air"`
	TaxEntities           string `gorm:"type:varchar(250);column:Tax_Entities" json:"Tax_Entities"`
	Stories               string `gorm:"type:varchar(250);column:Stories" json:"Stories"`
	TADTarrantAccountType string `gorm:"type:text;column:TAD_Tarrant_Account_Type" json:"TAD_Tarrant_Account_Type"`
	PropertyClass         string `gorm:"type:text;column:Property_Class" json:"Property_Class"`
	LandValue             int    `gorm:"type:int(13);column:Land_Value" json:"Land_Value"`
	ImprovementValue      int    `gorm:"type:int(11);column:Improvement_Value" json:"Improvement_Value"`
	TotalValue            int    `gorm:"type:int(11);column:Total_Value" json:"Total_Value"`
	AccountNum            int    `gorm:"type:int(11);column:Account_Num" json:"Account_Num"`
	OwnerName             string `gorm:"type:text;column:Owner_Name" json:"Owner_Name"`
	OwnerAddress          string `gorm:"type:text;column:Owner_Address" json:"Owner_Address"`
	OwnerCityState        string `gorm:"type:text;column:Owner_CityState" json:"Owner_CityState"`
	OwnerZip              int    `gorm:"type:int(5);column:Owner_Zip" json:"Owner_Zip"`
	OwnerZip4             int    `gorm:"type:int(4);column:Owner_Zip4" json:"Owner_Zip4"`
	TADMap                string `gorm:"type:text;column:TAD_Map" json:"TAD_Map"`
	Mapsco                string `gorm:"type:text;column:MAPSCO" json:"MAPSCO"`
	StateUseCode          string `gorm:"type:text;column:State_Use_Code" json:"State_Use_Code"`
	LegalDescription      string `gorm:"type:text;column:LegalDescription" json:"LegalDescription"`
	NumSpecialDist        int    `gorm:"type:int(1);column:Num_Special_Dist" json:"Num_Special_Dist"`
	// add sql.NullTime  mean value can nil
	DeedDate               sql.NullTime             `gorm:"type:date;column:Deed_Date;default:NULL" json:"Deed_Date"`
	AppraisalDate          sql.NullTime             `gorm:"type:date;column:Appraisal_Date;default:NULL" json:"Appraisal_Date"`
	AppraisedValue         int                      `gorm:"type:int(13);column:Appraised_Value" json:"Appraised_Value"`
	GISLink                string                   `gorm:"type:text;column:GIS_Link" json:"GIS_Link"`
	City                   string                   `gorm:"type:varchar(250);column:City" json:"City"`
	Units                  string                   `gorm:"type:varchar(250);column:Units" json:"Units"`
	CrawlingStatus         string                   `gorm:"type:varchar(50)" json:"crawling_status"`
	ZillowPriceHistory     []ZillowPriceHistory     `gorm:"foreignKey:Maindb3ID;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
	ZillowPublicTaxHistory []ZillowPublicTaxHistory `gorm:"foreignKey:Maindb3ID;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
}
