package extension

import (
	"fmt"
	"sample-service/app/common"
	esimDto "sample-service/app/esim/dto"
	esimRepo "sample-service/app/esim/repository"
	"sample-service/app/extension/dto"
	"sample-service/app/extension/repository"
	vendor "sample-service/app/internal_vendor"
	packagesDto "sample-service/app/packages/dto"
	packagesRepo "sample-service/app/packages/repository"

	"github.com/gofiber/fiber/v2"
)

func GetExtension(c *fiber.Ctx) error {

	println("Get Extension")

	// Parse request
	var requestQuery dto.GetExtensionRequestDTO
	if err := c.QueryParser(&requestQuery); err != nil {
		return common.ServerError(c, fmt.Sprintf("Error parsing request query: %v", err))
	}

	if err := common.ValidateStruct(&requestQuery); err != nil {
		return common.ServerError(c, fmt.Sprintf("Validation error: %v", err))
	}

	fmt.Printf("JSON query: %+v\n", requestQuery)

	isActualBalance := common.GetBool(requestQuery.ActualBalance)
	iccid := common.GetStr(requestQuery.Iccid)

	fmt.Printf("isActualBalance: %+v\n", isActualBalance)

	if isActualBalance && iccid == "" {
		return common.ServerError(c, fmt.Sprintf("For actual balance ICCID is needed"))
	}

	var response []dto.GetExtensionResponseDTO
	var extensions []dto.GetExtensionDbDTO
	var err error

	extensions, err = repository.GetExtension(requestQuery)

	if err != nil {
		return common.ServerError(c, fmt.Sprintf("DB while getting extensions: %v", err))
	}

	if extensions == nil {
		return c.JSON("Not found")
	}

	for _, e := range extensions {

		var status dto.ExtensionStatus

		if e.IsExpired {
			status = dto.ExtensionStatusNotActive
		} else {
			status = dto.ExtensionStatusActive
		}

		getPackagesRequest := packagesDto.GetPackageRequestDTO{
			Id: &e.PackageId,
		}

		packages, err := packagesRepo.GetPackagesWithAreaCount(getPackagesRequest)

		if err != nil {
			return common.ServerError(c, fmt.Sprintf("Error while getting GetPackagesWithAreaCount: %s", err.Error()))
		}

		dbPackage := packages[0]

		region, err := packagesRepo.GetRegionAndCountries(dbPackage.CountryList)

		if err != nil {
			continue
		}

		resp_ext := dto.GetExtensionResponseDTO{
			ID:                e.ID,
			PackageID:         e.PackageId,
			Iccid:             e.Iccid,
			Lpa:               e.Lpa,
			RegionCodeName:    region.CodeName,
			RegionName:        region.Name,
			Coverage:          int64(region.IncludedCountriesAmount),
			Days:              int64(dbPackage.Duration),
			ValueBytesStart:   e.InitialQuantityInBytes,
			ValueBytesCurrent: e.RemainingQuantityInBytes,
			DateStart:         e.StartTime,
			DateStop:          &e.EndTime,
			Status:            status,
		}

		response = append(response, resp_ext)
	}

	// Getting from Vendor actual balance, possible when icccid defined
	if isActualBalance && iccid != "" && len(response) == 1 {

		println("Get Extension Actual Balance")

		// getEsim form Vendor
		vendorQuery := esimDto.VendorGetEsimRequestDTO{
			Iccid: *requestQuery.Iccid,
		}
		vendorEsim, err := vendor.GetEsim(vendorQuery)

		if err != nil {
			return common.ServerError(c, fmt.Sprintf("Error getting actual balance: %v", err))
		}

		// We consider that plans from Vendor go in historical order
		for _, vendorPlan := range vendorEsim.AssignedPlans {
			// Save current balance to db
			repository.UpdateAssignedPlan(vendorPlan, vendorEsim.Iccid)
		}

		// Update ESIM status
		updEsim := esimDto.UpdateEsimRequestDTO{
			ProfileStatus: &vendorEsim.ProfileStatus,
		}

		_ = esimRepo.UpdateEsim(updEsim, iccid)
	}

	return c.JSON(response)

}
