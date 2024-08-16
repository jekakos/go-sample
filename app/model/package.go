package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	ID          int64            `gorm:"primaryKey;autoIncrement:true" json:"id"`
	VendorId    int64            `gorm:"unique" json:"vendor_id"` // 2435 Package ID
	GlobalId    string           `gorm:"unique" json:"global_id"` // "vendor-2435-30days-10gb-151"
	DataAmount  int64            `json:"data_amount"`
	Duration    int16            `json:"duration"`
	VendorPrice decimal.Decimal  `gorm:"type:numeric(10,2)" json:"vendor_price"` //
	CustomPrice *decimal.Decimal `gorm:"type:numeric(10,2)" json:"custom_price"` //
	Areas       []Area           `gorm:"many2many:package_areas;" json:"areas"`

	VendorCustomerPrice *decimal.Decimal `gorm:"type:numeric(10,2)" json:"vendor_customer_price"`
	Hidden              bool             `gorm:"default:false" json:"hidden"`
	Deleted             bool             `gorm:"default:false" json:"deleted"`
}
