package esim

import (
	"fmt"
	"sample-service/app/common"
	"sample-service/app/esim/dto"
	esimRepo "sample-service/app/esim/repository"
	vendor "sample-service/app/internal_vendor"
	packagesDto "sample-service/app/packages/dto"
	packagesRepo "sample-service/app/packages/repository"

	"github.com/gofiber/fiber/v2"
)

// @Summary Create ESIM
// @Description Create ESIM
// @Tags root
// @Accept application/json
// @Produce json
// @Param   request body dto.CreateEsimRequestDTO true "Create ESIM request payload"
// @Success 200     {object} dto.CreateEsimResponseDTO
// @Router /esim [post]
func CreateEsim(c *fiber.Ctx) error {

	println("Create ESIM")

	// Parse request
	var requestBody dto.CreateEsimRequestDTO
	if err := c.BodyParser(&requestBody); err != nil {
		return common.ServerError(c, fmt.Sprintf("Error parsing request body: %v", err))
	}

	if err := common.ValidateStruct(&requestBody); err != nil {
		return common.ServerError(c, fmt.Sprintf("Validation error: %v", err))
	}

	fmt.Printf("Vendor JSON body: %+v\n", requestBody)

	packageId, err := common.GetPackageFromGlobalId(requestBody.PackageId)

	if err != nil {
		return common.ServerError(c, fmt.Sprintf("Validation error: %v", err))
	}

	vendorId := common.GetVendorId()
	if requestBody.VendorId != vendorId {
		return common.ServerError(c, fmt.Sprintf("Invalid vendor ID: %d", requestBody.VendorId))
	}

	// Check package
	getPackageParams := packagesDto.GetPackageRequestDTO{
		Id: &requestBody.PackageId,
	}

	packages, err := packagesRepo.GetPackagesWithAreaCount(getPackageParams)
	if err != nil || len(packages) != 1 {
		return common.ServerError(c, fmt.Sprintf("Cannot find Package: %s", requestBody.PackageId))
	}

	// Create ESIM of Vendor side
	vendorRequest := dto.VendorCreateEsimRequestDTO{
		PlanId: packageId,
	}
	isProd := common.GetBool(requestBody.Prod)

	vendorResponse, err := vendor.CreateEsim(vendorRequest, isProd)
	if err != nil {
		return common.ServerError(c, fmt.Sprintf("Creation error: %v", err))
	}

	// Insert into db
	esimModel, err := esimRepo.AddEsim(requestBody.UserUUID, packageId, requestBody.PackageId, vendorResponse)
	if err != nil {
		return common.ServerError(c, fmt.Sprintf("E-Sim was created but wasn't saved in db  %v", err))
	}

	response := dto.CreateEsimResponseDTO{
		Iccid:       vendorResponse.Iccid,
		Lpa:         esimModel.Lpa,
		Imsi:        nil,
		Msisdn:      nil,
		ExtensionId: &esimModel.AssignedPlans[0].ID,
		UserUUID:    requestBody.UserUUID,
		VendorId:    requestBody.VendorId,
		PackageId:   requestBody.PackageId,
	}

	return c.JSON(response)
}
