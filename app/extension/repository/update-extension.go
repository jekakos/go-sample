package repository

import (
	"fmt"
	"sample-service/app/esim/dto"
	"sample-service/app/infra"
	"sample-service/app/model"
	"time"
)

func UpdateAssignedPlan(planDTO dto.AssignedPlansDTO, iccid string) error {

	println("Update extension")
	// Find Esim with AssignedPlans
	var esim model.Esim
	if err := infra.DB.Where("iccid = ?", iccid).First(&esim).Error; err != nil {
		return fmt.Errorf("failed to find Esim with ICCID %s: %w", iccid, err)
	}

	if len(esim.AssignedPlans) == 0 {
		return fmt.Errorf("no assigned plans found for Esim with iccid %s", iccid)
	}

	// Find Plan to update
	for _, dbPlan := range esim.AssignedPlans {
		if dbPlan.VendorPlanId == planDTO.Id {
			updatePlan := &dbPlan

			if planDTO.InitialQuantityInBytes != 0 {
				updatePlan.InitialQuantityInBytes = planDTO.InitialQuantityInBytes
			}
			if planDTO.RemainingQuantityInBytes != 0 {
				updatePlan.RemainingQuantityInBytes = planDTO.RemainingQuantityInBytes
			}
			if planDTO.StartTime != "" {
				startTime, err := time.Parse(time.RFC3339, planDTO.StartTime)
				if err != nil {
					return fmt.Errorf("invalid start time format: %w", err)
				}
				updatePlan.StartTime = &startTime
			}
			if planDTO.EndTime != "" {
				endTime, err := time.Parse(time.RFC3339, planDTO.EndTime)
				if err != nil {
					return fmt.Errorf("invalid end time format: %w", err)
				}
				updatePlan.EndTime = &endTime
			}
			updatePlan.IsExpired = planDTO.IsExpired

			// Update Ares
			if len(planDTO.Areas) > 0 {
				var areas []model.Area
				for _, areaDTO := range planDTO.Areas {
					area := model.Area{
						ID:     areaDTO.Id,
						Name:   areaDTO.Name,
						Region: areaDTO.Region,
						Iso:    &areaDTO.Iso,
					}
					areas = append(areas, area)
				}
				// Update links
				if err := infra.DB.Model(updatePlan).Association("Areas").Replace(areas); err != nil {
					return fmt.Errorf("failed to update areas: %w", err)
				}
			}

			// Save
			if err := infra.DB.Save(updatePlan).Error; err != nil {
				return fmt.Errorf("failed to update AssignedPlan: %w", err)
			}

		}
	}

	return nil
}
