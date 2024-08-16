package model

import (
	"time"

	"gorm.io/gorm"
)

type Esim struct {
	gorm.Model
	ID            int64              `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Iccid         string             `gorm:"unique;not null" json:"iccid"`
	Lpa           string             `gorm:"not null" json:"lpa"`
	SmdpAddress   string             `gorm:"not null" json:"smdpAddress"`
	ProfileStatus string             `gorm:"not null" json:"profileStatus"`
	UserUUID      string             `gorm:"not null" json:"user_uuid"`
	ProfileStaus  string             `json:"profileStaus,omitempty"`
	InstalledAt   *time.Time         `json:"installedAt,omitempty"`
	AssignedPlans []EsimAssignedPlan `gorm:"foreignKey:EsimID" json:"assignedPlans"`
}
