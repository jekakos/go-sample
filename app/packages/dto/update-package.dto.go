package dto

import "github.com/shopspring/decimal"

type UpdatePackageRequestDTO struct {
	VendorPrice     *decimal.Decimal `json:"vendor_price,omitempty"`
	Days            *int32           `json:"days,omitempty"`
	Hidden          *bool            `json:"hidden,omitempty"`
	CustomPrice     *decimal.Decimal `json:"custom_price,omitempty"`
	NullCustomPrice *bool            `json:"null_custom_price,omitempty"`
}

type UpdatePackageResponseDTO struct {
	Id                  string           `json:"id"`
	DataGb              decimal.Decimal  `json:"data_gb"`
	Days                int64            `json:"days"`
	VendorCustomerPrice *decimal.Decimal `json:"vendor_customer_price"`
	CustomPrice         *decimal.Decimal `json:"custom_price"`
	VendorPrice         decimal.Decimal  `json:"vendor_price"`
	Hidden              bool             `json:"hidden"`
	RegionName          string           `json:"region_name"`
	RegionCodeName      string           `json:"region_code_name"`
}