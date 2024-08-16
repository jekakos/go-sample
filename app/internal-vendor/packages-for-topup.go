package vendor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sample-service/app/internal_vendor/dto"
)

func GetPackagesForTopup(iccid string) ([]dto.VendorPackage, error) {

	var resultPackages []dto.VendorPackage

	token := os.Getenv("VENDOR_TOKEN")
	baseURL := os.Getenv("VENDOR_API_BASE")
	client := &http.Client{}

	url := fmt.Sprintf("%s/wholesale/esims/%s/top-up", baseURL, iccid)

	println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Apikey %s", token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Vendor error: unexpected status code: %d", resp.StatusCode)
	}

	var response dto.VendorPackageResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Decore JSON error: %s", err.Error()))
	}

	resultPackages = response.Rows

	return resultPackages, nil

}
