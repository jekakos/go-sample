package esim

import (
	"fmt"
	"sample-service/app/common"
	"sample-service/app/esim/dto"
	"sample-service/app/esim/repository"
	vendor "sample-service/app/internal_vendor"

	"github.com/gofiber/fiber/v2"
)

func GetEsim(c *fiber.Ctx) error {

	println("Get ESIM")

	// Parse request
	var requestQuery dto.GetEsimRequestDTO
	if err := c.QueryParser(&requestQuery); err != nil {
		return common.ServerError(c, fmt.Sprintf("Error parsing request query: %v", err))
	}

	if err := common.ValidateStruct(&requestQuery); err != nil {
		return common.ServerError(c, fmt.Sprintf("Validation error: %v", err))
	}
	fmt.Printf("JSON query: %+v\n", requestQuery)

	dbResponse, err := repository.GetEsim(requestQuery)

	if err != nil {
		return common.ServerError(c, fmt.Sprintf("DB while getting esims: %v", err))
	}

	if dbResponse == nil || len(dbResponse) == 0 {
		return c.JSON("Not found")
	}

	actualStatus := common.GetBool(requestQuery.ActualStatus)

	// Getting from Vendor actual balance, possible when icccid defined
	if actualStatus == true && *requestQuery.Iccid != "" && len(dbResponse) == 1 {

		println("Get Esim Actual Status")

		// getEsim form Vendor
		vendorQuery := dto.VendorGetEsimRequestDTO{
			Iccid: *requestQuery.Iccid,
		}
		vendorEsim, err := vendor.GetEsim(vendorQuery)

		if err != nil {
			return common.ServerError(c, fmt.Sprintf("Error getting actual status: %v", err))
		}

		updateData := dto.UpdateEsimRequestDTO{
			ProfileStatus: &vendorEsim.ProfileStatus,
			InstalledAt:   &vendorEsim.InstalledAt,
		}

		println(vendorEsim.ProfileStatus)
		// Save current status to db
		repository.UpdateEsim(updateData, *requestQuery.Iccid)

	}
	var response []dto.GetEsimResponseDTO
	for _, dbEsim := range dbResponse {
		esim := dto.GetEsimResponseDTO{
			Iccid:    dbEsim.Iccid,
			Lpa:      dbEsim.Lpa,
			Imsi:     "",
			Msisdn:   "",
			UserUUID: dbEsim.UserUUID,
			VendorId: common.GetVendorId(),
		}
		response = append(response, esim)
	}
	return c.JSON(response)

}
