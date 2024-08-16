package dto

import (
	"sample-service/app/esim/dto"
	"time"
)

type GetExtensionRequestDTO struct {
	ID            *int32  `json:"id,omitempty" query:"id"`
	Iccid         *string `json:"iccid,omitempty"  query:"iccid"`
	PackageID     *string `json:"package_id,omitempty"  query:"package_id"`
	UserUUID      *string `json:"user_uuid,omitempty"  query:"user_uuid"`
	ActualBalance *bool   `json:"actual_balance,omitempty"  query:"actual_balance"`
}

type ExtensionStatus string

const (
	ExtensionStatusActive    ExtensionStatus = "active"
	ExtensionStatusNotActive ExtensionStatus = "not_active"
	ExtensionStatusDisabled  ExtensionStatus = "disabled"
	ExtensionStatusDeleted   ExtensionStatus = "deleted"
	ExtensionStatusNone      ExtensionStatus = "none"
)

type GetExtensionResponseDTO struct {
	ID                int64           `json:"id"`
	Iccid             string          `json:"iccid"`
	Lpa               string          `json:"lpa"`
	RegionCodeName    string          `json:"region_code_name"`
	RegionName        string          `json:"region_name"`
	Coverage          int64           `json:"coverage"`
	PackageID         string          `json:"package_id"`
	Days              int64           `json:"days"`
	ValueBytesStart   int64           `json:"value_bytes_start"`
	ValueBytesCurrent int64           `json:"value_bytes_current"`
	DateStart         time.Time       `json:"date_start"`
	DateStop          *time.Time      `json:"date_stop,omitempty"`
	Status            ExtensionStatus `json:"status"`
	AddedAt           time.Time       `json:"added_at"`
}

type GetExtensionDbDTO struct {
	ID                       int64     `json:"id"`
	EsimID                   int64     `json:"esim_id"`
	PackageId                string    `json:"package_id"`
	InitialQuantityInBytes   int64     `json:"initial_quantity_in_bytes"`
	RemainingQuantityInBytes int64     `json:"remaining_quantity_in_bytes"`
	StartTime                time.Time `json:"start_time"`
	EndTime                  time.Time `json:"end_time"`
	IsExpired                bool      `json:"is_expired"`
	UserUUID                 string    `json:"user_uuid"`
	Iccid                    string    `json:"iccid"`
	Lpa                      string    `json:"lpa"`
}

type AddExtensionDbDTO struct {
	EsimID                   int64      `json:"esim_id"`
	PackageId                string     `json:"package_id"`
	InitialQuantityInBytes   int64      `json:"initial_quantity_in_bytes"`
	RemainingQuantityInBytes int64      `json:"remaining_quantity_in_bytes"`
	StartTime                *time.Time `json:"start_time"`
	EndTime                  *time.Time `json:"end_time"`
	IsExpired                bool       `json:"is_expired"`
	UserUUID                 string     `json:"user_uuid"`
	Iccid                    string     `json:"iccid"`
	VendorPlanId             int64      `json:"vendor_plan_id"`
	VendorAreaId             int64      `json:"vendor_area_id"`
	VendorId                 int64      `json:"vendor_id"`

	Areas []dto.AreaDTO `json:"areas"`
}

type AddExtensionRequestDTO struct {
	Iccid     string `json:"iccid" validate:"required"`
	PackageId string `json:"package_id" validate:"required"`
}
