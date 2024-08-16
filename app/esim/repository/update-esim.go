package repository

import (
	"fmt"
	"sample-service/app/common"
	"sample-service/app/esim/dto"
	"sample-service/app/infra"
	"sample-service/app/model"
	"time"
)

func UpdateEsim(esimData dto.UpdateEsimRequestDTO, iccid string) error {

	println("Save esim")
	// Find Esim
	var esim model.Esim
	if err := infra.DB.Where("iccid = ?", common.ClearQuery(iccid)).First(&esim).Error; err != nil {
		return fmt.Errorf("failed to find Esim with ICCID %s: %w", iccid, err)
	}
	fmt.Printf("ESIM: %+v\n", esim)

	updates := make(map[string]interface{})

	profileStatus := common.GetStr(esimData.ProfileStatus)
	fmt.Printf("Status = %s\n", profileStatus)

	if profileStatus != "" {
		updates["profile_status"] = profileStatus
	}

	if esimData.InstalledAt != nil {
		installedAt, err := time.Parse(time.RFC3339, *esimData.InstalledAt)
		if err == nil {
			updates["installed_at"] = installedAt
		}
	}

	fmt.Printf("Updates: %+v\n", updates)

	// Update
	if err := infra.DB.Model(&esim).Updates(updates).Error; err != nil {
		return fmt.Errorf("Failed to update esim: %w", err)
	}

	return nil
}
