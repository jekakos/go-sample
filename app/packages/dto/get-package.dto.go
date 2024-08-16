package dto

import "github.com/shopspring/decimal"

type GetPackageRequestDTO struct {
	Id             *string  `json:"id" query:"id"` // vendor-2605-7days-1gb where 2605 - package Vendor ID
	ShowHidden     *bool    `json:"show_hidden" query:"show_hidden"`
	RegionCodeName *string  `json:"region_code_name" query:"region_code_name"`
	RegionName     *string  `json:"region_name,omitempty" query:"region_name"`
	VendorIds      *[]int32 `json:"vendor_ids"`
}

type GetPackageResponseDTO struct {
	Id                  string           `json:"id"`
	DataGb              decimal.Decimal  `json:"data_gb"`
	Days                int64            `json:"days"`
	VendorCustomerPrice *decimal.Decimal `json:"vendor_customer_price"`
	CustomPrice         *decimal.Decimal `json:"custom_price"`
	VendorPrice         decimal.Decimal  `json:"vendor_price"`
	Hidden              bool             `json:"hidden"`
	Region              RegionDTO        `json:"region"`
}

type RegionDTO struct {
	CodeName                string       `json:"code_name"`
	Name                    string       `json:"name"`
	IncludedCountriesAmount int8         `json:"included_countries_amount"`
	Countries               []CountryDTO `json:"countries"`
}

type CountryDTO struct {
	CodeName string `json:"code_name"`
	Name     string `json:"name"`
}

/*
pub struct GetPackageRequestInterfaceDTO {
    pub id: Option<String>,
    pub show_hidden: Option<bool>,
    pub region_code_name: Option<String>,
}

pub struct GetPackageResponseInterfaceDTO {
    pub id: String,
    pub data_gb: Decimal, // i.e. 0.5 / 30.0
    pub days: i32,
    pub vendor_customer_price: Option<Decimal>,
    pub custom_price: Option<Decimal>,
    pub vendor_price: Decimal,
    pub hidden: bool,
    pub region: RegionInterfaceDTO,
}
pub struct RegionInterfaceDTO {
    pub code_name: String,
    pub name: String,
    pub included_countries_amount: i64,
 		pub countries: Option<Vec<CountryInterfaceDTO>>
}

pub struct CountryInterfaceDTO {
    pub code_name: String,
    pub name: String,
}
*/
