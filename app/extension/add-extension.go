package extension

import (
	"fmt"
	"sample-service/app/common"
	esimDto "sample-service/app/esim/dto"
	esimRepository "sample-service/app/esim/repository"
	extensionDto "sample-service/app/extension/dto"
	extRepository "sample-service/app/extension/repository"
	vendor "sample-service/app/internal_vendor"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get Packages for Top-up
// @Description Get Packages for Top-up
// @Tags root
// @Accept application/json
// @Produce json
// @Param   iccid path string true "Iccid"
// @Param   request body extensionDto.AddExtensionRequestDTO true "Add extesion request payload"
// @Success 200     {object} extensionDto.GetExtensionResponseDTO
// @Router /extension [post]
func AddExtensionTopup(c *fiber.Ctx) error {

	// Parse request
	var requestBody extensionDto.AddExtensionRequestDTO
	if err := c.BodyParser(&requestBody); err != nil {
		return common.ServerError(c, fmt.Sprintf("Error parsing request body: %v", err))
	}

	if err := common.ValidateStruct(&requestBody); err != nil {
		return common.ServerError(c, fmt.Sprintf("Validation error: %v", err))
	}

	iccid := requestBody.Iccid

	if !common.IsValidICCID(iccid) {
		return common.ServerError(c, fmt.Sprintf("Invalid ICCID format: %s", iccid))
	}

	planId, err := common.GetPackageFromGlobalId(requestBody.PackageId)

	if err != nil {
		return common.ServerError(c, fmt.Sprintf("Invalid Package ID format: %s", requestBody.PackageId))
	}

	// Get ESIm from db
	esimParams := esimDto.GetEsimRequestDTO{
		Iccid: &iccid,
	}

	dBesims, err := esimRepository.GetEsim(esimParams)

	if err != nil || len(dBesims) != 1 {
		return common.ServerError(c, fmt.Sprintf("Connot find E-Sim by iccid: %s", iccid))
	}

	dbEsim := dBesims[0]
	fmt.Printf("ESIM: %+v\n", dbEsim)

	// Add plan (ext) on Vendor side
	requestVendor := esimDto.VendorCreateEsimRequestDTO{
		PlanId: planId,
	}

	esimWithNewPlan, err := vendor.EsimTopup(iccid, requestVendor)
	if err != nil {
		return common.ServerError(c, fmt.Sprintf("Error while EsimTopup on Vendor side: %s", err.Error()))
	}

	// Suppose that new added plan should be the first one
	newPlan := esimWithNewPlan.AssignedPlans[0]

	fmt.Printf("New plan: %+v\n", newPlan)

	// Add to db AssignedPlan
	addPlanDbDto := extensionDto.AddExtensionDbDTO{
		EsimID:                   dbEsim.ID,
		PackageId:                requestBody.PackageId,
		InitialQuantityInBytes:   newPlan.InitialQuantityInBytes,
		RemainingQuantityInBytes: newPlan.RemainingQuantityInBytes,
		StartTime:                common.ParseTime(newPlan.StartTime),
		EndTime:                  common.ParseTime(newPlan.EndTime),
		IsExpired:                newPlan.IsExpired,
		UserUUID:                 dbEsim.UserUUID,
		Iccid:                    dbEsim.Iccid,
		Areas:                    newPlan.Areas,
		VendorPlanId:             newPlan.PlanId,
		VendorId:                 newPlan.Id,
	}

	newPannDb, err := extRepository.AddExtension(addPlanDbDto)
	if err != nil {
		return common.ServerError(c, fmt.Sprintf("Error while AddExtensionDbDTO: %s", err.Error()))
	}

	var status extensionDto.ExtensionStatus

	if newPannDb.IsExpired {
		status = extensionDto.ExtensionStatusNotActive
	} else {
		status = extensionDto.ExtensionStatusActive
	}

	newPlanResponse := extensionDto.GetExtensionResponseDTO{
		ID:                newPannDb.ID,
		PackageID:         requestBody.PackageId,
		Iccid:             iccid,
		ValueBytesStart:   newPannDb.InitialQuantityInBytes,
		ValueBytesCurrent: newPannDb.RemainingQuantityInBytes,
		DateStart:         common.GetTime(newPannDb.StartTime),
		DateStop:          newPannDb.EndTime,
		Status:            status,
	}

	c.Set("Content-Type", "application/json")
	return c.JSON(newPlanResponse)
}
