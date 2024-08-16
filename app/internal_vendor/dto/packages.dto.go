package dto

import "github.com/shopspring/decimal"

type VendorPackage struct {
	Id         int64           `json:"id"`
	DataAmount int64           `json:"dataAmount"`
	Duration   int16           `json:"duration"`
	Price      decimal.Decimal `json:"price"`
	Areas      []VendorArea    `json:"areas"`
}

type VendorArea struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
	Iso    string `json:"iso"`
}

type VendorPackageResponse struct {
	Count int             `json:"count"`
	Rows  []VendorPackage `json:"rows"`
}
