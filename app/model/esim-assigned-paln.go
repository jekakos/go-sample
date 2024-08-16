package model

import (
	"time"

	"gorm.io/gorm"
)

type EsimAssignedPlan struct {
	gorm.Model
	ID                       int64      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	VendorId                 int64      `json:"vendor_id"`
	VendorPlanId             int64      `json:"vendor_plan_id"`
	VendorAreaId             int64      `json:"vendor_area_id"`
	EsimID                   int64      `gorm:"not null" json:"esim_id"`
	PackageId                string     `json:"package_id"`
	InitialQuantityInBytes   int64      `gorm:"not null" json:"initialQuantityInBytes"`
	RemainingQuantityInBytes int64      `gorm:"not null" json:"remainingQuantityInBytes"`
	StartTime                *time.Time `json:"startTime,omitempty"`
	EndTime                  *time.Time `json:"endTime,omitempty"`
	IsExpired                bool       `gorm:"not null" json:"isExpired"`
	Areas                    []Area     `gorm:"many2many:esim_assigned_plan_areas" json:"areas"`
}
