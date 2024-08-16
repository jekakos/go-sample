package model

import (
	"gorm.io/gorm"
)

type Area struct {
	gorm.Model
	ID       int64              `gorm:"primaryKey" json:"id"` //This is Vendor ID, not autoincrement
	Name     string             `gorm:"index:idx_name_region_iso,unique" json:"name"`
	Region   string             `gorm:"index:idx_name_region_iso,unique" json:"region"`
	Iso      *string            `gorm:"index:idx_name_region_iso,unique" json:"iso"`
	Packages []Package          `gorm:"many2many:package_areas;" json:"packages"`
	Plans    []EsimAssignedPlan `gorm:"many2many:esim_assigned_plan_areas;" json:"plans"`
}
